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
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/tendermint/tendermint/types"
	"sync"
)

var (
	// how many block each goroutine need to sync when do fast sync
	syncBlockNumFastSync = int64(conf.SyncBlockNumFastSync)

	// limit max goroutine
	limitChan = make(chan int64, conf.SyncMaxGoroutine)

	mutex           sync.Mutex
	mutexWatchBlock sync.Mutex
)

// start sync server
func Start() {
	var (
		status *ctypes.ResultStatus
		err    error
		i      = 1
	)
	Init()
	c := helper.GetClient().Client

	for {
		logger.Info.Printf("Begin %v time fast sync task", i)
		syncLatestHeight := fastSync(c)
		status, err = c.Status()
		if err != nil {
			logger.Error.Printf("TmClient err and try again, %v\n", err.Error())
			c := helper.GetClient().Client
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

	startCron(c)
}

func Init() {
	store.InitWithAuth()

	chainId := conf.ChainId
	syncTask, err := document.QuerySyncTask()

	if err != nil {
		if chainId == "" {
			logger.Error.Fatalln("sync process start failed,chainId is empty")
		}
		syncTask = document.SyncTask{
			Height:  0,
			ChainID: chainId,
		}
		store.Save(syncTask)
	}

	// init client pool
	helper.InitClientPool()
}

// start cron scheduler
func startCron(client rpcClient.Client) {
	spec := conf.SyncCron
	c := cron.New()
	c.AddFunc(spec, func() {
		watchBlock(client)
	})
	go c.Start()
}

func watchBlock(c rpcClient.Client) {
	mutexWatchBlock.Lock()

	syncTask, _ := document.QuerySyncTask()
	status, _ := c.Status()
	latestBlockHeight := status.SyncInfo.LatestBlockHeight

	// note: interval two block, to avoid get can't delegation at latest block
	//       sdk of this version may has some problem
	if syncTask.Height+2 <= latestBlockHeight {
		funcChain := []func(tx store.Docs, mutex sync.Mutex){
			handler.SaveTx, handler.SaveAccount, handler.UpdateBalance,
		}

		ch := make(chan int64)
		limitChan <- 1

		go syncBlock(syncTask.Height+1, latestBlockHeight-1, funcChain, ch, 0, constant.SyncTypeWatch)

		syncedLatestBlockHeight := latestBlockHeight - 1
		block, _ := c.Block(&syncedLatestBlockHeight)
		syncTask.Height = block.Block.Height
		syncTask.Time = block.Block.Time

		select {
		case <-ch:
			logger.Info.Printf("Watch block, current height is %v \n", latestBlockHeight)
			err := store.Update(syncTask)
			if err != nil {
				logger.Error.Printf("Update syncTask fail, err is %v",
					err.Error())
			}
		}
	} else {
		logger.Info.Printf("%v: wait, synced height is %v, latest height is %v\n",
			constant.SyncTypeWatch, syncTask.Height, latestBlockHeight)
	}

	mutexWatchBlock.Unlock()
}

// fast sync data from blockChain
func fastSync(c rpcClient.Client) int64 {
	syncTaskDoc, _ := document.QuerySyncTask()
	status, _ := c.Status()
	latestBlockHeight := status.SyncInfo.LatestBlockHeight

	funcChain := []func(tx store.Docs, mutex sync.Mutex){
		handler.SaveTx, handler.SaveAccount, handler.UpdateBalance,
	}

	ch := make(chan int64)

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
		go syncBlock(start, end, funcChain, ch, i, constant.SyncTypeFastSync)
	}

	for {
		select {
		case threadNo := <-ch:
			activeGoroutineNum = activeGoroutineNum - 1
			logger.Info.Printf("ThreadNo[%d] is over and active thread num is %d\n", threadNo, activeGoroutineNum)
			if activeGoroutineNum == 0 {
				goto end
			}
		}
	}

end:
	{
		logger.Info.Println("This fastSync task complete!")
		// update syncTask document
		block, _ := c.Block(&latestBlockHeight)
		syncTaskDoc.Height = block.Block.Height
		syncTaskDoc.Time = block.Block.Time
		err := store.Update(syncTaskDoc)
		if err != nil {
			logger.Error.Printf("Update syncTask fail, err is %v",
				err.Error())
		}
		return syncTaskDoc.Height
	}
}

func syncBlock(start int64, end int64, funcChain []func(tx store.Docs, mutex sync.Mutex),
	ch chan int64, threadNum int64, syncType string) {
	logger.Info.Printf("%v: ThreadNo[%d] begin sync block from %d to %d\n",
		syncType, threadNum, start, end)

	client := helper.GetClient()
	// release client
	defer client.Release()

	for j := start; j <= end; j++ {
		block, err := client.Client.Block(&j)
		if err != nil {
			logger.Error.Printf("Invalid block height %d and err is %v, try again\n", j, err.Error())
			// try again
			client2 := helper.GetClient()
			block, err = client2.Client.Block(&j)
			if err != nil {
				ch <- threadNum
				logger.Error.Fatalf("Invalid block height %d and err is %v\n", j, err.Error())
			}
		}
		if block.BlockMeta.Header.NumTxs > 0 {
			txs := block.Block.Data.Txs
			for _, txByte := range txs {
				docTx := helper.ParseTx(codec.Cdc, txByte, block.Block)
				txHash := helper.BuildHex(txByte.Hash())
				if txHash == "" {
					logger.Warning.Printf("Tx has no hash, skip this tx."+
						""+"tx is %v\n", helper.ToJson(docTx))
					continue
				}
				logger.Info.Printf("===========threadNo[%d] find tx, txHash=%s\n", threadNum, txHash)

				handler.Handle(docTx, mutex, funcChain)
			}
		}

		// get validatorSet at given height
		var validators []*types.Validator
		res, err := client.Client.Validators(&j)
		if err != nil {
			logger.Error.Printf("Can't get validatorSet at %v\n", j)
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
		syncType, threadNum, start, end)

	<-limitChan
	ch <- threadNum
	logger.Info.Printf("Send threadNum into channel: %v\n", threadNum)

}
