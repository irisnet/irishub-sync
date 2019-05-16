package task

import (
	"fmt"
	"os"
	"time"

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
)

func StartExecuteTask() {
	var (
		syncConfModel           document.SyncConf
		blockNumPerWorkerHandle int64
		maxWorkerSleepTime      int64
	)
	log := logger.GetLogger("TaskExecutor")

	// get sync conf
	syncConf, err := syncConfModel.GetConf()
	if err != nil {
		log.Fatal("Get sync conf failed", logger.String("err", err.Error()))
	}
	blockNumPerWorkerHandle = syncConf.BlockNumPerWorkerHandle
	if blockNumPerWorkerHandle <= 0 {
		log.Fatal("blockNumPerWorkerHandle should greater than 0")
	}
	maxWorkerSleepTime = syncConf.MaxWorkerSleepTime
	if maxWorkerSleepTime <= 0 {
		log.Fatal("maxWorkerSleepTime should greater than 0")
	}

	log.Info("Start execute task", logger.Any("sync conf", syncConf))

	// buffer channel to limit goroutine num
	chanLimit := make(chan bool, serverConf.WorkerNumExecuteTask)

	for {
		time.Sleep(time.Duration(1) * time.Second)
		chanLimit <- true
		go executeTask(blockNumPerWorkerHandle, maxWorkerSleepTime, chanLimit)
	}
}

func executeTask(blockNumPerWorkerHandle, maxWorkerSleepTime int64, chanLimit chan bool) {
	var (
		syncTaskModel          document.SyncTask
		workerId, taskType     string
		blockChainLatestHeight int64
	)
	log := logger.GetLogger("TaskExecutor")
	genWorkerId := func() string {
		// generate worker id use hostname@xxx
		hostname, _ := os.Hostname()
		return fmt.Sprintf("%v@%v", hostname, bson.NewObjectId().Hex())
	}

	healthCheckQuit := make(chan bool)
	workerId = genWorkerId()
	client := helper.GetClient()

	defer func() {
		if r := recover(); r != nil {
			log.Error("execute task fail", logger.Any("err", r))
		}
		close(healthCheckQuit)
		<-chanLimit
		client.Release()
	}()

	// check whether exist executable task
	// status = unhandled or
	// status = underway and now - lastUpdateTime > confTime
	tasks, err := syncTaskModel.GetExecutableTask(maxWorkerSleepTime)
	if err != nil {
		log.Error("Get executable task fail", logger.String("err", err.Error()))
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
			log.Info("Task has been take over by other goroutine")
		} else {
			log.Error("Take over task fail", logger.String("err", err.Error()))
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
	log.Info("worker begin execute task",
		logger.String("cur_worker", workerId), logger.String("task_id", task.ID.Hex()),
		logger.String("from-to", fmt.Sprintf("%v-%v", task.StartHeight, task.EndHeight)))

	// worker health check, if worker is alive, then update last update time every minute.
	// health check will exit in follow conditions:
	// 1. task is not owned by current worker
	// 2. task is invalid
	workerHealthCheck := func(taskId bson.ObjectId, currentWorker string) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("worker health check err", logger.Any("err", r))
			}
		}()

		for {
			select {
			case <-healthCheckQuit:
				logger.Info("get health check quit signal, now exit health check")
				return
			default:
				task, err := syncTaskModel.GetTaskByIdAndWorker(taskId, workerId)
				if err == nil {
					blockChainLatestHeight, err := getBlockChainLatestHeight()
					if err == nil {
						if assertTaskValid(task, blockNumPerWorkerHandle, blockChainLatestHeight) {
							// update task last update time
							if err := syncTaskModel.UpdateLastUpdateTime(task); err != nil {
								log.Error("update last update time fail", logger.String("err", err.Error()),
									logger.String("task_id", task.ID.Hex()))
							}
							logger.Info("health check success, now sleep one minute",
								logger.String("task_id", task.ID.Hex()),
								logger.String("task_current_worker", task.WorkerId))
						} else {
							log.Info("task is invalid, exit health check", logger.String("task_id", taskId.Hex()))
							break
						}
					} else {
						log.Error("get block chain latest height fail", logger.String("err", err.Error()))
					}
				} else {
					if err == mgo.ErrNotFound {
						log.Info("task may be task over by other goroutine, exit health check",
							logger.String("task_id", taskId.Hex()), logger.String("current_worker", workerId))
						break
					} else {
						log.Error("get task by id and worker fail", logger.String("task_id", taskId.Hex()),
							logger.String("current_worker", workerId))
					}
				}
				time.Sleep(1 * time.Minute)
			}
		}
	}
	go workerHealthCheck(task.ID, workerId)

	// check task is valid
	// valid catch up task: current_height < end_height
	// valid follow task: current_height + blockNumPerWorkerHandle > blockChainLatestHeight
	blockChainLatestHeight, err = getBlockChainLatestHeight()
	if err != nil {
		log.Error("get block chain latest height fail", logger.String("err", err.Error()))
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
				log.Error("get block chain latest height fail", logger.String("err", err.Error()))
				return
			}

			if task.CurrentHeight+2 >= blockChainLatestHeight {
				// wait block chain latest block height updated, must interval two block
				log.Info("wait block chain latest block height updated, must interval two block",
					logger.String("taskId", task.ID.String()),
					logger.String("workerId", task.WorkerId),
					logger.Int64("taskCurrentHeight", task.CurrentHeight),
					logger.Int64("blockChainLatestHeight", blockChainLatestHeight))
				time.Sleep(2 * time.Second)
				continue
			}
		}

		// parse block and tx
		blockDoc, err := parseBlock(inProcessBlock, client)
		if err != nil {
			log.Error("Parse block fail", logger.Int64("block", inProcessBlock),
				logger.String("err", err.Error()))
		}

		// check task owner
		workerUnchanged, err := assertTaskWorkerUnchanged(task.ID, task.WorkerId)
		if err != nil {
			log.Error("assert task worker is unchanged fail", logger.String("err", err.Error()))
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
				log.Error("save docs fail", logger.String("err", err.Error()))
			} else {
				task.CurrentHeight = inProcessBlock

				if taskType == document.SyncTaskTypeFollow {
					// TODO: whether can remove compareAndUpdateValidators in sync logic
					// compare and update validators
					handler.CompareAndUpdateValidators()
				}
			}
		} else {
			log.Info("task worker changed", logger.Any("task_id", task.ID),
				logger.String("origin worker", workerId), logger.String("current worker", task.WorkerId))
			return
		}
	}

	log.Info("worker finish execute task",
		logger.String("task_worker", task.WorkerId), logger.Any("task_id", task.ID),
		logger.String("from-to-current", fmt.Sprintf("%v-%v-%v", task.StartHeight, task.EndHeight, task.CurrentHeight)))
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
	var blockDoc document.Block

	defer func() {
		if err := recover(); err != nil {
			logger.Error("parse block fail", logger.Int64("blockHeight", b),
				logger.Any("err", err))
		}
	}()

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

	accsBalanceNeedUpdatedByParseTxs, err := handler.HandleTx(block.Block)
	if err != nil {
		return blockDoc, err
	}

	// get validatorSet at given height
	var validators []*types.Validator
	res, err := client.Validators(&b)
	if err != nil {
		logger.Error("Can't get validatorSet at height", logger.Int64("Height", b))
	} else {
		validators = res.Validators
	}

	return handler.ParseBlock(block.BlockMeta, block.Block, validators, accsBalanceNeedUpdatedByParseTxs), nil
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