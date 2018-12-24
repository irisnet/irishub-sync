package task

import (
	"fmt"
	serverConf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/service/handler"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/helper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"os"
	"sync"
	"time"
)

func StartExecuteTask() {
	var (
		syncConfModel           document.SyncConf
		blockNumPerWorkerHandle int64
		maxWorkerSleepTime      int64
	)
	// get sync conf
	syncConf, err := syncConfModel.GetConf()
	if err != nil {
		logger.Fatal("Get sync conf failed", logger.String("err", err.Error()))
	}
	blockNumPerWorkerHandle = syncConf.BlockNumPerWorkerHandle
	if blockNumPerWorkerHandle <= 0 {
		logger.Fatal("blockNumPerWorkerHandle should greater than 0")
	}
	maxWorkerSleepTime = syncConf.MaxWorkerSleepTime
	if maxWorkerSleepTime <= 0 {
		logger.Fatal("maxWorkerSleepTime should greater than 0")
	}

	logger.Info("Start execute task", logger.Any("sync conf", syncConf))

	// buffer channel to limit goroutine num
	chanLimit := make(chan bool, serverConf.WorkerNumExecuteTask)

	for {
		chanLimit <- true
		go executeTask(blockNumPerWorkerHandle, maxWorkerSleepTime, chanLimit)
	}
}

func executeTask(blockNumPerWorkerHandle, maxWorkerSleepTime int64, chanLimit chan bool) {
	var (
		syncTaskModel          document.SyncTask
		workerId               string
		taskType               string
		blockChainLatestHeight int64
	)

	genWorkerId := func() string {
		// generate worker id use hostname@xxx
		hostname, _ := os.Hostname()
		return fmt.Sprintf("%v@%v", hostname, bson.NewObjectId().Hex())
	}

	workerId = genWorkerId()
	client := helper.GetClient()

	defer func() {
		if r := recover(); r != nil {
			logger.Error("execute task fail", logger.Any("err", r))
		}
		client.Release()
		<-chanLimit
	}()

	// check sync task if exist
	// status = unhandled or
	// status = underway and now - lastUpdateTime > confTime
	tasks, err := syncTaskModel.GetExecutableTask(maxWorkerSleepTime)
	if err != nil {
		logger.Error("Get executable task fail", logger.String("err", err.Error()))
	}
	if len(tasks) == 0 {
		// there is no executable tasks
		return
	}

	// take over sync task
	// attempt to update status, worker_id and worker_logs
	task := tasks[0]
	err = syncTaskModel.TakeOverTask(task, workerId)
	if err != nil {
		if err == mgo.ErrNotFound {
			// this task has been take over by other goroutine
			logger.Info("Task has been take over by other goroutine")
		} else {
			logger.Error("Take over task fail", logger.String("err", err.Error()))
		}
		return
	} else {
		// task over task success, update task worker to current worker
		task.WorkerId = workerId
	}

	if task.EndHeight != 0 {
		taskType = document.SyncTaskTypeCatchUp
	} else {
		taskType = document.SyncTaskTypeFollow
	}
	logger.Info("worker begin execute task", logger.String("cur_worker", workerId),
		logger.String("task_type", taskType), logger.Any("task", task))

	// check task is valid
	// valid catch up task: current_height < end_height
	// valid follow task: current_height + blockNumPerWorkerHandle > blockChainLatestHeight
	blockChainLatestHeight, err = getBlockChainLatestHeight()
	if err != nil {
		logger.Error("get block chain latest height fail", logger.String("err", err.Error()))
		return
	}
	for assertTaskValid(task, blockNumPerWorkerHandle, blockChainLatestHeight) {
		var inProcessBlock int64
		if task.CurrentHeight == 0 {
			inProcessBlock = task.StartHeight
		} else {
			inProcessBlock = task.CurrentHeight + 1
		}

		// if task is follow task,
		// wait value of blockChainLatestHeight updated when inProcessBlock >= blockChainLatestHeight
		if taskType == document.SyncTaskTypeFollow {
			blockChainLatestHeight, err = getBlockChainLatestHeight()
			if err != nil {
				logger.Error("get block chain latest height fail", logger.String("err", err.Error()))
				return
			}

			if task.CurrentHeight+2 >= blockChainLatestHeight {
				// wait block chain latest block height updated, must interval two block
				continue
			}
		}

		// parse block and tx
		blockDoc, err := parseBlock(inProcessBlock, client)
		if err != nil {
			logger.Error("Parse block fail", logger.Int64("block", inProcessBlock),
				logger.String("err", err.Error()))
		}

		// check task owner
		workerUnchanged, err := assertTaskWorkerUnchanged(task.ID, task.WorkerId)
		if err != nil {
			logger.Error("assert task worker is unchanged fail", logger.String("err", err.Error()))
		}
		if workerUnchanged {
			// save data and update sync task
			taskDoc := task
			taskDoc.CurrentHeight = inProcessBlock
			taskDoc.LastUpdateTime = time.Now().Unix()
			taskDoc.Status = document.SyncTaskStatusUnderway
			if inProcessBlock == task.EndHeight {
				taskDoc.Status = document.SyncTaskStatusCompleted
			}

			err := saveDocs(blockDoc, taskDoc)
			if err != nil {
				logger.Error("save docs fail", logger.String("err", err.Error()))
			} else {
				task.CurrentHeight = inProcessBlock

				if taskType == document.SyncTaskTypeFollow {
					// compare and update validators
					handler.CompareAndUpdateValidators()
				}
			}
		} else {
			logger.Info("task worker changed", logger.String("origin worker", workerId),
				logger.String("current worker", task.WorkerId))
			return
		}
	}

	logger.Info("worker finish execute task", logger.String("task_worker", task.WorkerId),
		logger.String("task_type", taskType), logger.Any("task", task.ID))
}

// assert task is valid
// valid catch up task: current_height < end_height
// valid follow task: current_height + blockNumPerWorkerHandle > blockChainLatestHeight
func assertTaskValid(task document.SyncTask, blockNumPerWorkerHandle, blockChainLatestHeight int64) bool {
	var (
		taskType string
		flag     = false
	)
	if task.EndHeight != 0 {
		taskType = document.SyncTaskTypeCatchUp
	} else {
		taskType = document.SyncTaskTypeFollow
	}
	currentHeight := task.CurrentHeight
	if currentHeight == 0 {
		currentHeight = task.StartHeight - 1
	}

	switch taskType {
	case document.SyncTaskTypeCatchUp:
		if currentHeight < task.EndHeight {
			flag = true
		}
		break
	case document.SyncTaskTypeFollow:
		if currentHeight+blockNumPerWorkerHandle > blockChainLatestHeight {
			flag = true
		}
		break
	}
	return flag
}

func parseBlock(b int64, client *helper.Client) (document.Block, error) {
	var (
		mutex    sync.Mutex
		blockDoc document.Block
	)

	// define functions which should be executed
	// during parse tx and block
	funcChain := []handler.Action{
		handler.SaveTx, handler.SaveAccount, handler.SaveOrUpdateDelegator,
	}

	block, err := client.Block(&b)
	if err != nil {
		// there is possible parse block fail when in iterator
		var err2 error
		client2 := helper.GetClient()
		block, err2 = client2.Block(&b)
		client2.Release()
		if err2 != nil {
			return blockDoc, err2
		}
	}

	// save or update common_tx, tx_msg, proposal, delegator, candidate, account document
	// TODO: saveOrUpdate above documents, save block and update sync task should be in a transaction.
	// TODO  this task will be finished during second refactor plan.
	if block.BlockMeta.Header.NumTxs > 0 {
		txs := block.Block.Data.Txs
		for _, txByte := range txs {
			docTx := helper.ParseTx(txByte, block.Block)
			txHash := helper.BuildHex(txByte.Hash())
			if txHash == "" {
				logger.Warn("Tx has no hash, skip this tx.", logger.Any("Tx", docTx))
				continue
			}
			handler.Handle(docTx, mutex, funcChain)
		}
	}

	// get validatorSet at given height
	var validators []*types.Validator
	res, err := client.Validators(&b)
	if err != nil {
		logger.Error("Can't get validatorSet at height", logger.Int64("Height", b))
	} else {
		validators = res.Validators
	}

	return handler.ParseBlock(block.BlockMeta, block.Block, validators), nil
}

// assert task worker unchanged
func assertTaskWorkerUnchanged(taskId bson.ObjectId, workerId string) (bool, error) {
	var (
		syncTaskModel document.SyncTask
	)
	// check task owner
	task, err := syncTaskModel.GetTaskById(taskId)
	if err != nil {
		return false, err
	}

	if task.WorkerId == workerId {
		return true, nil
	} else {
		return false, nil
	}
}

func saveDocs(blockDoc document.Block, taskDoc document.SyncTask) error {
	var (
		ops []txn.Op
	)

	if blockDoc.Hash == "" {
		return fmt.Errorf("block document is empty")
	}

	insertOp := txn.Op{
		C:      document.CollectionNmBlock,
		Id:     bson.NewObjectId(),
		Insert: blockDoc,
	}

	updateOp := txn.Op{
		C:      document.CollectionNameSyncTask,
		Id:     taskDoc.ID,
		Assert: txn.DocExists,
		Update: bson.M{
			"$set": bson.M{
				"current_height":   taskDoc.CurrentHeight,
				"status":           taskDoc.Status,
				"last_update_time": taskDoc.LastUpdateTime,
			},
		},
	}

	ops = append(ops, insertOp, updateOp)

	if len(ops) > 0 {
		err := store.Txn(ops)
		if err != nil {
			return err
		}
	}

	return nil
}
