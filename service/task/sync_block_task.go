package task

import (
	"fmt"
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/service/handler"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
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
	locker    *sync.Mutex
	HasTask   bool
	StartTime time.Time
}

func (watcher *BlockWatcher) Lock() {
	watcher.locker.Lock()
	watcher.HasTask = true
	watcher.StartTime = time.Now()
}

func (watcher *BlockWatcher) UnLock() {
	watcher.HasTask = false
	watcher.locker.Unlock()
}

func (watcher *BlockWatcher) CanDo() bool {
	if !watcher.HasTask {
		return true
	}
	logger.Debug("========================task's trigger [watchBlock] hashTask===================")
	sub := time.Now().Sub(watcher.StartTime).Seconds()
	if sub < 10 {
		logger.Debug("========================task's trigger [watchBlock] skip Task===================")
		return false
	}
	logger.Debug("======================== task's trigger [watchBlock] hashTask,reset HasTask flag ========================")
	watcher.UnLock()
	return true
}

func NewWatcher() *BlockWatcher {
	return watcher
}

func MakeSyncBlockTask() Task {
	return NewLockTaskFromEnv(conf.CronWatchBlock, "watch_block_lock_key_lock", func() {
		if watcher.CanDo() {
			logger.Debug("========================task's trigger [watchBlock] begin===================")
			watcher.watchBlock()
			logger.Debug("========================task's trigger [watchBlock] end===================")
		}
	})
}

// start watcher
func (watcher *BlockWatcher) FastSync() {
	var (
		status *types.ResultStatus
		err    error
		i      = 1
	)
	client := helper.GetClient()
	defer client.Release()

	// fast sync
	for {
		logger.Info("Begin fast sync task", logger.Int("time", i))
		syncLatestHeight := watcher.fastSync()
		status, err = client.Status()
		if err != nil {
			logger.Panic("TmClient err and try again, %v\n", logger.String("err", err.Error()))
		}
		latestHeight := status.SyncInfo.LatestBlockHeight
		if syncLatestHeight >= latestHeight-60 {
			logger.Info("All fast sync task complete!")
			break
		}
		logger.Info("End fast sync task", logger.Int("time", i))
		i++
	}
}

// start cron scheduler

func (watcher *BlockWatcher) watchBlock() {
	methodName = constant.SyncTypeWatch
	watcher.Lock()
	client := helper.GetClient()

	defer func() {
		if err := recover(); err != nil {
			logger.Error("task watchBlock execute error", logger.Any("err", err))
		}
		logger.Debug("debug=======================5 watchBlock defer=======================debug")
		client.Release()
		logger.Debug("debug=======================6 watchBlock defer client.Release()=======================debug")
		watcher.UnLock()
		logger.Debug("debug=======================7 watchBlock defer watcher.UnLock()=======================debug")
	}()

	logger.Debug("debug=======================1=======================debug")

	status, err := client.Status()
	if err != nil {
		logger.Error("TmClient err and try again, %v\n", logger.String("err", err.Error()))
		return
	}
	logger.Debug("debug=======================2=======================debug")
	syncTask, _ := document.QuerySyncTask()
	logger.Debug("debug=======================3=======================debug")
	latestBlockHeight := status.SyncInfo.LatestBlockHeight

	// note: interval two block, to avoid get can't delegation at latest block
	//       sdk of this version may has some problem
	if syncTask.Height+1 <= latestBlockHeight-1 {
		logger.Info("latest height", logger.String("method", methodName), logger.Int64("latestBlockHeight", latestBlockHeight))
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
			logger.Info("latest height", logger.String("method", constant.SyncTypeWatch), logger.Int64("syncedLatestBlockHeight", syncedLatestBlockHeight))
			if err := store.Update(syncTask); err != nil {
				logger.Error("Update syncTask fail", logger.String("method", constant.SyncTypeWatch), logger.String("err", err.Error()))
			}
		}
	} else {
		logger.Info("system's speed of syncing is fast than blockchain,must be interval two block", logger.Int64("syncHeight", syncTask.Height), logger.Int64("latestHeight", latestBlockHeight))
	}
	logger.Debug("debug=======================4=======================debug")
}

// fast sync data from blockChain
func (watcher *BlockWatcher) fastSync() int64 {
	methodName = constant.SyncTypeFastSync

	client := helper.GetClient()
	defer client.Release()

	status, err := client.Status()
	if err != nil {
		logger.Error("TmClient err", logger.String("err", err.Error()))
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
			logger.Info("server is syncing,will try again next 10 second")
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
			logger.Info("Thread is over", logger.Int64("threadNo", threadNo), logger.Int64("activeGoroutineNum", activeGoroutineNum))
			if activeGoroutineNum == 0 {
				goto end
			}
		}
	}

end:
	{
		logger.Info("This fastSync task complete!")
		// update syncTask document
		block, _ := client.Block(&latestBlockHeight)
		syncTaskDoc.Height = block.Block.Height
		syncTaskDoc.Time = block.Block.Time
		syncTaskDoc.Syncing = false
		err := store.Update(syncTaskDoc)
		if err != nil {
			logger.Error("Update syncTask fail", logger.String("err", err.Error()))
		}
		return syncTaskDoc.Height
	}
}

func syncBlock(start, end, threadNum int64,
	ch chan int64, syncType string,
	funcChain []handler.Action) {

	methodName = fmt.Sprintf("syncBlock_%s", syncType)

	logger.Info("begin sync block", logger.Int64("threadNo", threadNum), logger.Int64("start", start), logger.Int64("end", end))
	client := helper.GetClient()

	defer func() {
		if err := recover(); err != nil {
			logger.Error("syncBlock err", logger.Any("err", err))
		}
		<-limitChan
		ch <- threadNum
		client.Release()
	}()

	for b := start; b <= end; b++ {
		block, err := client.Block(&b)
		if err != nil {
			logger.Error("acquire block error", logger.Int64("height", b), logger.String("err", err.Error()))
			// try again
			client2 := helper.GetClient()
			block, err = client2.Block(&b)
			if err != nil {
				ch <- threadNum
				logger.Error("acquire block error", logger.Int64("height", b), logger.String("err", err.Error()))
			}
		}
		if block.BlockMeta.Header.NumTxs > 0 {
			txs := block.Block.Data.Txs
			for _, txByte := range txs {
				docTx := helper.ParseTx(codec.Cdc, txByte, block.Block)
				txHash := helper.BuildHex(txByte.Hash())
				if txHash == "" {
					logger.Warn("Tx has no hash, skip this tx.", logger.Any("Tx", docTx))
					continue
				}
				logger.Info("found tx", logger.Int64("threadNo", threadNum), logger.String("hash", txHash))
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

		// save block info
		logger.Info("save block", logger.Int64("threadNo", threadNum), logger.Int64("Height", block.Block.Height))
		handler.SaveBlock(block.BlockMeta, block.Block, validators)

		// compare and update validators during watch block
		if syncType == constant.SyncTypeWatch {
			handler.CompareAndUpdateValidators(validators)
		}
	}

	logger.Info("finish sync block", logger.Int64("threadNo", threadNum), logger.Int64("from", start), logger.Int64("to", end))

	logger.Info("Send threadNum into channel", logger.Int64("threadNo", threadNum))
}
