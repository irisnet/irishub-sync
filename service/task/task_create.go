package task

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

func StartCreateTask() {
	log := logger.GetLogger("StartCreateTask")
	var (
		syncConfModel           document.SyncConf
		blockNumPerWorkerHandle int64
	)

	// get sync conf
	syncConf, err := syncConfModel.GetConf()
	if err != nil {
		log.Fatal("Get sync conf failed", logger.String("err", err.Error()))
	}
	blockNumPerWorkerHandle = syncConf.BlockNumPerWorkerHandle
	if blockNumPerWorkerHandle <= 0 {
		log.Fatal("blockNumPerWorkerHandle should greater than 0")
	}

	log.Info("Start create task", logger.Any("sync conf", syncConf))

	// buffer channel to limit goroutine num
	chanLimit := make(chan bool, serverConf.WorkerNumCreateTask)

	for {
		time.Sleep(time.Duration(1) * time.Second)
		chanLimit <- true
		go createTask(blockNumPerWorkerHandle, chanLimit)
	}
}

func createTask(blockNumPerWorker int64, chanLimit chan bool) {
	var (
		syncTaskModel     document.SyncTask
		syncTasks         []document.SyncTask
		ops               []txn.Op
		invalidFollowTask document.SyncTask
	)
	log := logger.GetLogger("createTask")

	defer func() {
		if err := recover(); err != nil {
			log.Error("Create sync task failed", logger.Any("err", err))
		}
		<-chanLimit
	}()

	// check valid follow task if exist
	// status of valid follow task is unhandled or underway
	validFollowTasks, err := syncTaskModel.QueryAll(
		[]string{
			document.SyncTaskStatusUnHandled,
			document.SyncTaskStatusUnderway,
		}, document.SyncTaskTypeFollow)
	if err != nil {
		log.Error("Query sync task failed", logger.String("err", err.Error()))
	}
	if len(validFollowTasks) == 0 {
		// get max end_height from sync_task
		maxEndHeight, err := syncTaskModel.GetMaxBlockHeight()
		if err != nil {
			log.Error("Get max end_block failed", logger.String("err", err.Error()))
			return
		}

		currentBlockHeight, err := getBlockChainLatestHeight()
		if err != nil {
			log.Error("Get current block height failed", logger.String("err", err.Error()))
			return
		}

		if maxEndHeight+blockNumPerWorker <= currentBlockHeight {
			syncTasks = createCatchUpTask(maxEndHeight, blockNumPerWorker, currentBlockHeight)
			log.Info("Create catch up task during follow task not exist", logger.Int64("from", maxEndHeight+1), logger.Int64("to", currentBlockHeight))
		} else {
			finished, err := assertAllCatchUpTaskFinished()
			if err != nil {
				log.Error("AssertAllCatchUpTaskFinished failed", logger.String("err", err.Error()))
				return
			}
			if finished {
				syncTasks = createFollowTask(maxEndHeight, blockNumPerWorker, currentBlockHeight)
				log.Info("Create follow task during follow task not exist", logger.Int64("from", maxEndHeight+1), logger.Int64("to", currentBlockHeight))
			}
		}
	} else {
		followTask := validFollowTasks[0]
		followedHeight := followTask.CurrentHeight
		if followedHeight == 0 {
			followedHeight = followTask.StartHeight - 1
		}

		currentBlockHeight, err := getBlockChainLatestHeight()
		if err != nil {
			log.Error("Get current block height failed", logger.String("err", err.Error()))
			return
		}

		if followedHeight+blockNumPerWorker <= currentBlockHeight {
			syncTasks = createCatchUpTask(followedHeight, blockNumPerWorker, currentBlockHeight)

			invalidFollowTask = followTask
			log.Info("Create catch up task during follow task exist",
				logger.Int64("from", followedHeight+1),
				logger.Int64("to", currentBlockHeight),
				logger.String("invalidFollowTask_id", invalidFollowTask.ID.Hex()),
				logger.Int64("invalidFollowTask_currentHeight", invalidFollowTask.CurrentHeight))
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

	if invalidFollowTask.ID.Valid() {
		op := txn.Op{
			C:  document.CollectionNameSyncTask,
			Id: invalidFollowTask.ID,
			Assert: bson.M{
				"current_height": invalidFollowTask.CurrentHeight,
			},
			Update: bson.M{
				"$set": bson.M{
					"status": document.FollowTaskStatusInvalid,
				},
			},
		}
		ops = append(ops, op)
	}

	if len(ops) > 0 {
		err := store.Txn(ops)
		if err != nil {
			log.Warn("Create sync task fail", logger.String("err", err.Error()))
		} else {
			log.Info("Create sync task success")
		}
	}
}

// get current block height
func getBlockChainLatestHeight() (int64, error) {
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

func createCatchUpTask(maxEndHeight, blockNumPerWorker, currentBlockHeight int64) []document.SyncTask {
	var (
		syncTasks []document.SyncTask
	)

	for maxEndHeight+blockNumPerWorker <= currentBlockHeight {
		syncTask := document.SyncTask{
			StartHeight:    maxEndHeight + 1,
			EndHeight:      maxEndHeight + blockNumPerWorker,
			Status:         document.SyncTaskStatusUnHandled,
			LastUpdateTime: time.Now().Unix(),
		}
		syncTasks = append(syncTasks, syncTask)

		maxEndHeight += blockNumPerWorker
	}

	return syncTasks
}

func assertAllCatchUpTaskFinished() (bool, error) {
	var (
		syncTaskModel          document.SyncTask
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
			LastUpdateTime: time.Now().Unix(),
		}

		syncTasks = append(syncTasks, syncTask)
	}

	return syncTasks
}
