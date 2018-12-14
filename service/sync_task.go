package service

import (
	serverConf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/helper"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"time"
)

const (
	goroutineNumCreateTask = 2
)

var (
	syncTaskModel           document.SyncTask
	syncConfModel           document.SyncConf
	blockNumPerWorkerHandle int64
)

func Start() {
	// get sync conf
	syncConf, err := syncConfModel.GetConf(serverConf.ChainId)
	if err != nil {
		logger.Fatal("Get sync conf failed", logger.String("err", err.Error()))
	}
	blockNumPerWorkerHandle = syncConf.BlockNumPerWorkerHandle
	if blockNumPerWorkerHandle <= 0 {
		logger.Fatal("blockNumPerWorkerHandle should greater than 0")
	}

	// buffer channel to limit goroutine num
	chanLimit := make(chan bool, goroutineNumCreateTask)

	for {
		chanLimit <- true
		go createTask(blockNumPerWorkerHandle, chanLimit)
	}
}

func createTask(blockNumPerWorker int64, chanLimit chan bool) {
	var (
		syncTasks  []document.SyncTask
		ops        []txn.Op
		removeTask document.SyncTask
	)

	defer func() {
		if err := recover(); err != nil {
			logger.Error("Create sync task failed", logger.Any("err", err))
		}
		<-chanLimit
	}()

	// get current block height
	getCurrentBlockHeight := func() (int64, error) {
		client := helper.GetClient()
		defer func() {
			client.Release()
		}()
		status, err := client.Status()
		if err != nil {
			return 0, err
		}
		currentBlockHeight := status.SyncInfo.LatestBlockHeight

		return currentBlockHeight, nil
	}

	// check follow task if exist
	followTasks, err := syncTaskModel.QueryAll([]string{}, document.SyncTaskTypeFollow)
	if err != nil {
		logger.Error("Query sync task failed", logger.String("err", err.Error()))
	}
	if len(followTasks) == 0 {
		// get max end_height from sync_task
		maxEndHeight, err := syncTaskModel.GetMaxBlockHeight()
		if err != nil {
			logger.Error("Get max end_block failed", logger.String("err", err.Error()))
			return
		}

		currentBlockHeight, err := getCurrentBlockHeight()
		if err != nil {
			logger.Error("Get current block height failed", logger.String("err", err.Error()))
			return
		}

		if maxEndHeight+blockNumPerWorker <= currentBlockHeight {
			syncTasks = createCatchUpTask(maxEndHeight, blockNumPerWorker, currentBlockHeight)
			logger.Info("Create catch up task during follow task not exist", logger.Int64("from", maxEndHeight), logger.Int64("to", currentBlockHeight))
		} else {
			finished, err := assertAllCatchUpTaskFinished()
			if err != nil {
				logger.Error("AssertAllCatchUpTaskFinished failed", logger.String("err", err.Error()))
				return
			}
			if finished {
				syncTasks = createFollowTask(maxEndHeight, blockNumPerWorker, currentBlockHeight)
				logger.Info("Create follow task during follow task not exist", logger.Int64("from", maxEndHeight), logger.Int64("to", currentBlockHeight))
			}
		}
	} else {
		followTask := followTasks[0]
		followedHeight := followTask.CurrentHeight

		currentBlockHeight, err := getCurrentBlockHeight()
		if err != nil {
			logger.Error("Get current block height failed", logger.String("err", err.Error()))
			return
		}

		if followedHeight+blockNumPerWorker <= currentBlockHeight {
			syncTasks = createCatchUpTask(followedHeight, blockNumPerWorker, currentBlockHeight)

			removeTask = followTask
			logger.Info("Create catch up task during follow task exist", logger.Int64("from", followedHeight), logger.Int64("to", currentBlockHeight))
		}
	}

	// bulk insert or remove use transaction
	if len(syncTasks) > 0 {
		for _, v := range syncTasks {
			objectId := bson.NewObjectId()
			v.ID = objectId
			op := txn.Op{
				C:      document.CollectionNameSyncTask,
				Id:     objectId,
				Assert: nil,
				Insert: v,
			}

			ops = append(ops, op)
		}
	}

	if removeTask.ID.Valid() {
		removeOp := txn.Op{
			C:      document.CollectionNameSyncTask,
			Id:     removeTask.ID,
			Assert: txn.DocExists,
			Remove: true,
		}
		ops = append(ops, removeOp)
	}

	if len(ops) > 0 {
		c := store.GetCollection(store.CollectionNameTxn)
		runner := txn.NewRunner(c)

		txObjectId := bson.NewObjectId()
		err := runner.Run(ops, txObjectId, nil)
		if err != nil {
			if err == txn.ErrAborted {
				err = runner.Resume(txObjectId)
				if err != nil {
					logger.Error("Resume transaction failed", logger.String("err", err.Error()))
				}
			} else {
				logger.Error("Unknown while run create sync task transaction", logger.String("err", err.Error()))

			}
		} else {
			logger.Info("create sync task success")
		}
	}
}

func createCatchUpTask(maxEndHeight, blockNumPerWorker, currentBlockHeight int64) []document.SyncTask {
	var (
		syncTasks []document.SyncTask
	)

	for maxEndHeight+blockNumPerWorker <= currentBlockHeight {
		syncTask := document.SyncTask{
			StartHeight:    maxEndHeight + 1,
			EndHeight:      maxEndHeight + blockNumPerWorker,
			Status:         document.SyncTaskStatusUnHandled,
			LastUpdateTime: time.Now(),
		}
		syncTasks = append(syncTasks, syncTask)

		maxEndHeight += blockNumPerWorker
	}

	return syncTasks
}

func assertAllCatchUpTaskFinished() (bool, error) {
	var (
		allCatchUpTaskFinished = false
	)

	// assert all catch up task whether finished
	tasks, err := syncTaskModel.QueryAll(
		[]string{
			document.SyncTaskStatusUnHandled,
			document.SyncTaskStatusUnderway,
		},
		document.SyncTaskTypeCatchUp)
	if err != nil {
		return false, err
	}

	if len(tasks) == 0 {
		allCatchUpTaskFinished = true
	}

	return allCatchUpTaskFinished, nil
}

func createFollowTask(maxEndHeight, blockNumPerWorker, currentBlockHeight int64) []document.SyncTask {
	var (
		syncTasks []document.SyncTask
	)

	if maxEndHeight+blockNumPerWorker > currentBlockHeight {
		syncTask := document.SyncTask{
			StartHeight:    maxEndHeight + 1,
			EndHeight:      0,
			Status:         document.SyncTaskStatusUnHandled,
			LastUpdateTime: time.Now(),
		}

		syncTasks = append(syncTasks, syncTask)
	}

	return syncTasks
}
