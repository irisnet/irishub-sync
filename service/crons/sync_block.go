package crons

import (
	"fmt"
	"github.com/irisnet/irishub-sync/service/handler"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
)

func watchBlock() {
	methodName = constant.SyncTypeWatch

	client := helper.GetClient()
	mutexWatchBlock.Lock()
	defer func() {
		mutexWatchBlock.Unlock()
		if err := recover(); err != nil {
			logger.Error.Println(err)
		}
		client.Release()
	}()
	status, err := client.Status()
	if err != nil {
		fmt.Println(err)
		return
	}

	syncTask, _ := document.QuerySyncTask()
	latestBlockHeight := status.SyncInfo.LatestBlockHeight

	// note: interval two block, to avoid get can't delegation at latest block
	//       sdk of this version may has some problem
	if syncTask.Height+1 <= latestBlockHeight-1 {
		logger.Info.Printf("%v: latest height is %v\n",
			methodName, latestBlockHeight)

		funcChain := []handler.Action{
			handler.SaveTx, handler.SaveAccount, handler.UpdateBalance,
		}

		ch := make(chan int64)
		limitChan <- 1

		go syncBlock(syncTask.Height+1, latestBlockHeight-1, 0, ch, constant.SyncTypeWatch, funcChain)

		syncedLatestBlockHeight := latestBlockHeight - 1
		block, _ := client.Block(&syncedLatestBlockHeight)
		syncTask.Height = syncedLatestBlockHeight
		syncTask.Time = block.Block.Time

		select {
		case <-ch:
			logger.Info.Printf("%v: synced height is %v \n",
				constant.SyncTypeWatch, syncedLatestBlockHeight)
			if err := store.Update(syncTask); err != nil {
				logger.Error.Printf("%v: Update syncTask fail, err is %v\n",
					methodName, err.Error())
			}
		}
	} else {
		logger.Info.Printf("%v: wait, synced height is %v, latest height is %v\n",
			methodName, syncTask.Height, latestBlockHeight)
	}

}
