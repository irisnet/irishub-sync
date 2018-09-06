package service

import (
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/service/handler"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/helper"

	"github.com/irisnet/irishub-sync/store/document"
	"github.com/robfig/cron"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"fmt"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/tendermint/tendermint/types"
	"sync"
)

var (
	// how many block each goroutine need to sync when do fast sync
	syncBlockNumFastSync = int64(conf.SyncBlockNumFastSync)

	// limit max goroutine
	limitChan = make(chan int64, conf.SyncMaxGoroutine)

	mutex, mutexWatchBlock sync.Mutex

	methodName string
)

func init() {
	// init store
	store.InitWithAuth()

	// init client pool
	chainId := conf.ChainId
	syncTask, err := document.QuerySyncTask()
	if err != nil {
		if chainId == "" {
			logger.Error.Fatalln("sync process start failed, chainId is empty")
		}
		syncTask = document.SyncTask{
			Height:  0,
			ChainID: chainId,
		}
		store.Save(syncTask)
	}

	helper.InitClientPool()
}

// start sync server
func Start() {
	var (
		status *ctypes.ResultStatus
		err    error
		i      = 1
	)
	client := helper.GetClient()
	defer client.Release()

	c := client.Client

	// fast sync
	for {
		logger.Info.Printf("Begin %v time fast sync task", i)
		syncLatestHeight := fastSync()
		status, err = c.Status()
		if err != nil {
			logger.Error.Printf("TmClient err and try again, %v\n", err.Error())
			client := helper.GetClient()
			c := client.Client
			status, err = c.Status()
			if err != nil {
				logger.Error.Fatalf("TmClient err and exit, %v\n", err.Error())
			}
		}
		latestHeight := status.SyncInfo.LatestBlockHeight
		if syncLatestHeight >= latestHeight-60 {
			logger.Info.Println("All fast sync task complete!")
			break
		}
		logger.Info.Printf("End %v time fast sync task", i)
		i++
	}

	// watch sync
	startCron()
}

// start cron scheduler
func startCron() {
	c := cron.New()
	c.AddFunc(conf.CronWatchBlock, func() {
		watchBlock()
	})
	c.AddFunc(conf.CronCalculateUpTime, func() {
		handler.CalculateAndSaveValidatorUpTime()
	})
	c.AddFunc(conf.CronCalculateTxGas, func() {
		handler.CalculateTxGasAndGasPrice()
	})
	c.AddFunc(conf.SyncProposalStatus, func() {
		handler.SyncProposalStatus()
	})
	go c.Start()
}

func watchBlock() {
	methodName = constant.SyncTypeWatch

	client := helper.GetClient()
	defer client.Release()

	mutexWatchBlock.Lock()

	c := client.Client
	status, _ := c.Status()

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
		block, _ := c.Block(&syncedLatestBlockHeight)
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

	mutexWatchBlock.Unlock()
}

// fast sync data from blockChain
func fastSync() int64 {
	methodName = constant.SyncTypeFastSync

	client := helper.GetClient()
	defer client.Release()

	c := client.Client
	status, err := c.Status()
	if err != nil {
		logger.Error.Printf("TmClient err, %v\n", err)
		return 0
	}

	// define functions which should be executed
	// during parse tx and block
	funcChain := []func(tx document.CommonTx, mutex sync.Mutex){
		handler.SaveTx, handler.SaveAccount, handler.UpdateBalance,
	}

	// define unbuffered channel
	ch := make(chan int64)

	// define how many goroutine should be used during fast sync
	syncTaskDoc, _ := document.QuerySyncTask()
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
		block, _ := c.Block(&latestBlockHeight)
		syncTaskDoc.Height = block.Block.Height
		syncTaskDoc.Time = block.Block.Time
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
	defer client.Release()

	for b := start; b <= end; b++ {
		block, err := client.Client.Block(&b)
		if err != nil {
			logger.Error.Printf("%v: Invalid block height %d and err is %v, try again\n",
				methodName, b, err.Error())
			// try again
			client2 := helper.GetClient()
			block, err = client2.Client.Block(&b)
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
		res, err := client.Client.Validators(&b)
		if err != nil {
			logger.Error.Printf("%v: Can't get validatorSet at %v\n", methodName, b)
		} else {
			validators = res.Validators
		}

		// save block info
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
