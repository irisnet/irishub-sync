package task

import (
	"fmt"
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/service/handler"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"sync"
	"time"
)

var (
	// how many block each goroutine need to sync when do fast sync
	syncBlockNumFastSync = int64(conf.SyncBlockNumFastSync)
	// limit max goroutine
	limitChan  = make(chan int64, conf.SyncMaxGoroutine)
	mutex      sync.Mutex
	methodName string

	watcher = &BlockWatcher{
		locker:  new(sync.Mutex),
		HasTask: false,
	}
)

type BlockWatcher struct {
	locker  *sync.Mutex
	HasTask bool
}

func (watcher *BlockWatcher) Lock() {
	watcher.locker.Lock()
	watcher.HasTask = true
	watcher.locker.Unlock()
}

func (watcher *BlockWatcher) UnLock() {
	watcher.locker.Lock()
	watcher.HasTask = false
	watcher.locker.Unlock()
}

func NewWatcher() *BlockWatcher {
	return watcher
}

func MakeSyncBlockTask() Task {
	return NewLockTaskFromEnv(conf.CronWatchBlock, "watch_block_lock_key_lock", func() {
		if watcher.HasTask {
			logger.Info.Printf("========================task's trigger [%s] hashTask===================", "watchBlock")
			return
		}
		logger.Info.Printf("========================task's trigger [%s] begin===================", "watchBlock")
		watcher.watchBlock()
		logger.Info.Printf("========================task's trigger [%s] end===================", "watchBlock")
	})
}

// start watcher
func (watcher *BlockWatcher) FastSync() {
	var (
		status *ctypes.ResultStatus
		err    error
		i      = 1
	)
	client := helper.GetClient()
	defer client.Release()

	// fast sync
	for {
		logger.Info.Printf("Begin %v time fast sync task", i)
		syncLatestHeight := watcher.fastSync()
		status, err = client.Status()
		if err != nil {
			logger.Error.Printf("TmClient err and try again, %v\n", err.Error())
			panic(err)
		}
		latestHeight := status.SyncInfo.LatestBlockHeight
		if syncLatestHeight >= latestHeight-60 {
			logger.Info.Println("All fast sync task complete!")
			break
		}
		logger.Info.Printf("End %v time fast sync task", i)
		i++
	}
}

// start cron scheduler

func (watcher *BlockWatcher) watchBlock() {
	methodName = constant.SyncTypeWatch
	watcher.Lock()
	client := helper.GetClient()

	defer func() {
		logger.Info.Println("debug=======================5 watchBlock defer=======================debug")
		client.Release()
		logger.Info.Println("debug=======================6 watchBlock defer client.Release()=======================debug")
		watcher.UnLock()
		logger.Info.Println("debug=======================7 watchBlock defer watcher.UnLock()=======================debug")
	}()

	logger.Info.Println("debug=======================1=======================debug")

	status, err := client.Status()
	if err != nil {
		logger.Error.Println(err)
		return
	}
	logger.Info.Println("debug=======================2=======================debug")
	syncTask, _ := document.QuerySyncTask()
	logger.Info.Println("debug=======================3=======================debug")
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
	logger.Info.Println("debug=======================4=======================debug")
}

// fast sync data from blockChain
func (watcher *BlockWatcher) fastSync() int64 {
	methodName = constant.SyncTypeFastSync

	client := helper.GetClient()
	defer client.Release()

	status, err := client.Status()
	if err != nil {
		logger.Error.Printf("TmClient err, %v\n", err)
		return 0
	}

	// define functions which should be executed
	// during parse tx and block
	funcChain := []handler.Action{
		handler.SaveTx, handler.SaveAccount, handler.UpdateBalance,
	}

	// define unbuffered channel
	ch := make(chan int64)
loop:
	// define how many goroutine should be used during fast sync
	syncTaskDoc, err := document.QuerySyncTask()
	if err != nil {
		syncTaskDoc = document.SyncTask{
			Height:  0,
			ChainID: conf.ChainId,
			Syncing: true,
		}
		store.Save(syncTaskDoc)
	} else {
		if syncTaskDoc.Syncing {
			logger.Info.Println("server is syncing,will try again next 10 second")
			time.Sleep(10 * time.Second)
			goto loop
		}
	}
	latestBlockHeight := status.SyncInfo.LatestBlockHeight

	goroutineNum := (latestBlockHeight - syncTaskDoc.Height) / syncBlockNumFastSync

	if goroutineNum == 0 {
		goroutineNum = 20
		syncBlockNumFastSync = (latestBlockHeight - syncTaskDoc.Height) / goroutineNum
	}
	activeGoroutineNum := goroutineNum

	for i := int64(1); i <= goroutineNum; i++ {
		limitChan <- i
		var (
			start = syncTaskDoc.Height + (i-1)*syncBlockNumFastSync + 1
			end   = syncTaskDoc.Height + i*syncBlockNumFastSync
		)
		if i == goroutineNum {
			end = latestBlockHeight
		}
		go syncBlock(start, end, i, ch, constant.SyncTypeFastSync, funcChain)
	}

	for {
		select {
		case threadNo := <-ch:
			activeGoroutineNum = activeGoroutineNum - 1
			logger.Info.Printf("%v: ThreadNo[%d] is over and active thread num is %d\n",
				methodName, threadNo, activeGoroutineNum)
			if activeGoroutineNum == 0 {
				goto end
			}
		}
	}

end:
	{
		logger.Info.Printf("%v: This fastSync task complete!", methodName)
		// update syncTask document
		block, _ := client.Block(&latestBlockHeight)
		syncTaskDoc.Height = block.Block.Height
		syncTaskDoc.Time = block.Block.Time
		syncTaskDoc.Syncing = false
		err := store.Update(syncTaskDoc)
		if err != nil {
			logger.Error.Printf("%v: Update syncTask fail, err is %v",
				methodName, err.Error())
		}
		return syncTaskDoc.Height
	}
}

func syncBlock(start, end, threadNum int64,
	ch chan int64, syncType string,
	funcChain []handler.Action) {

	methodName = fmt.Sprintf("syncBlock_%s", syncType)

	logger.Info.Printf("%v: ThreadNo[%d] begin sync block from %d to %d\n",
		methodName, threadNum, start, end)

	client := helper.GetClient()

	defer func() {
		if err := recover(); err != nil {
			logger.Error.Println(err)
		}
		client.Release()
	}()

	for b := start; b <= end; b++ {
		block, err := client.Block(&b)
		if err != nil {
			logger.Error.Printf("%v: Invalid block height %d and err is %v, try again\n",
				methodName, b, err.Error())
			// try again
			client2 := helper.GetClient()
			block, err = client2.Block(&b)
			if err != nil {
				ch <- threadNum
				logger.Error.Fatalf("%v: Invalid block height %d and err is %v\n",
					methodName, b, err.Error())
			}
		}
		if block.BlockMeta.Header.NumTxs > 0 {
			txs := block.Block.Data.Txs
			for _, txByte := range txs {
				docTx := helper.ParseTx(codec.Cdc, txByte, block.Block)
				txHash := helper.BuildHex(txByte.Hash())
				if txHash == "" {
					logger.Warning.Printf("%v: Tx has no hash, skip this tx."+
						""+"tx is %v\n", methodName, helper.ToJson(docTx))
					continue
				}
				logger.Info.Printf("%v: ====ThreadNo[%d] find tx, txHash=%s\n",
					methodName, threadNum, txHash)

				handler.Handle(docTx, mutex, funcChain)
			}
		}

		// get validatorSet at given height
		var validators []*types.Validator
		res, err := client.Validators(&b)
		if err != nil {
			logger.Error.Printf("%v: Can't get validatorSet at %v\n", methodName, b)
		} else {
			validators = res.Validators
		}

		// save block info
		logger.Info.Printf("thread[%d] save block,height:[%d]", threadNum, block.Block.Height)
		handler.SaveBlock(block.BlockMeta, block.Block, validators)

		// compare and update validators during watch block
		if syncType == constant.SyncTypeWatch {
			handler.CompareAndUpdateValidators(validators)
		}
	}

	logger.Info.Printf("%v: ThreadNo[%d] finish sync block from %d to %d\n",
		methodName, threadNum, start, end)

	<-limitChan
	ch <- threadNum
	logger.Info.Printf("%v: Send threadNum into channel: %v\n",
		methodName, threadNum)
}
