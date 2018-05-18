package sync

import (
	"encoding/hex"
	"strings"

	conf "github.com/irisnet/iris-sync-server/conf/server"
	"github.com/irisnet/iris-sync-server/model/store"
	"github.com/irisnet/iris-sync-server/module/logger"
	"github.com/irisnet/iris-sync-server/module/stake"
	"github.com/irisnet/iris-sync-server/util/constant"
	"github.com/irisnet/iris-sync-server/util/helper"

	"github.com/irisnet/iris-sync-server/model/store/document"
	"github.com/robfig/cron"
	rpcClient "github.com/tendermint/tendermint/rpc/client"

	"sync"
)

var (
	// how many block each goroutine need to sync when do fast sync
	syncBlockNumFastSync = int64(conf.SyncBlockNumFastSync)

	// limit max goroutine
	// limitChan = make(chan int64, conf.SyncMaxGoroutine)

	mutex sync.Mutex
)

// start sync server
func Start() {
	InitServer()
	c := helper.GetClient().Client
	if err := fastSync(c); err != nil {
		logger.Error.Fatalf("sync block failed,%v\n", err)
	}
	startCron(c)
}

func InitServer() {
	store.Init()

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

func watchBlock(c rpcClient.Client) error {
	syncTask, _ := document.QuerySyncTask()
	status, _ := c.Status()
	latestBlockHeight := status.LatestBlockHeight

	// for test
	// latestBlockHeight := int64(60010)

	funcChain := []func(tx store.Docs, mutex sync.Mutex){
		saveTx, saveOrUpdateAccount, updateAccountBalance,
	}

	ch := make(chan int64)
	// limitChan <- 1

	go syncBlock(syncTask.Height+1, latestBlockHeight, funcChain, ch, 0)

	select {
	case <-ch:
		logger.Info.Printf("Watch block, current height is %v \n", latestBlockHeight)
		block, _ := c.Block(&latestBlockHeight)
		syncTask.Height = block.Block.Height
		syncTask.Time = block.Block.Time
		return store.Update(syncTask)
	}
}

// fast sync data from blockChain
func fastSync(c rpcClient.Client) error {
	syncTaskDoc, _ := document.QuerySyncTask()
	status, _ := c.Status()
	latestBlockHeight := status.LatestBlockHeight

	// for test
	// latestBlockHeight := int64(60000)

	funcChain := []func(tx store.Docs, mutex sync.Mutex){
		saveTx, saveOrUpdateAccount, updateAccountBalance,
	}

	ch := make(chan int64)
	activeThreadNum := int64(0)

	goRoutineNum := (latestBlockHeight - syncTaskDoc.Height) / syncBlockNumFastSync

	if goRoutineNum == 0 {
		goRoutineNum = 10
		syncBlockNumFastSync = 100
	}

	for i := int64(1); i <= goRoutineNum; i++ {
		activeThreadNum++
		// limitChan <- i
		var (
			start = syncTaskDoc.Height + (i-1)*syncBlockNumFastSync + 1
			end   = syncTaskDoc.Height + i*syncBlockNumFastSync
		)
		if i == goRoutineNum {
			end = latestBlockHeight
		}
		go syncBlock(start, end, funcChain, ch, i)
	}

	for {
		select {
		case threadNo := <-ch:
			activeThreadNum = activeThreadNum - 1
			logger.Info.Printf("ThreadNo[%d] is over and active thread num is %d\n", threadNo, activeThreadNum)
			if activeThreadNum == 0 {
				goto end
			}
		}
	}

end:
	{
		logger.Info.Println("Fast sync block, complete sync task")
		// update syncTask document
		block, _ := c.Block(&latestBlockHeight)
		syncTaskDoc.Height = block.Block.Height
		syncTaskDoc.Time = block.Block.Time
		store.Update(syncTaskDoc)
		return nil
	}
}

func syncBlock(start int64, end int64, funcChain []func(tx store.Docs, mutex sync.Mutex), ch chan int64, threadNum int64) {
	logger.Info.Printf("ThreadNo[%d] begin sync block from %d to %d\n",
		threadNum, start, end)
	
	client := helper.GetClient()
	// release client
	defer client.Release()
	// release unBuffer chain and buffer chain
	defer func() {
		ch <- threadNum
		logger.Info.Printf("Send threadNum into channel: %v\n", threadNum)
		// <- limitChan
	}()

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
				txType, tx := helper.ParseTx(txByte)
				txHash := strings.ToUpper(hex.EncodeToString(txByte.Hash()))
				if txHash == "" {
					logger.Warning.Printf("Tx has no hash, skip this tx."+
						""+"type of tx is %v, value of tx is %v\n",
						txType, tx)
					continue
				}
				logger.Info.Printf("===========threadNo[%d] find tx,txType=%s;txHash=%s\n", threadNum, txType, txHash)

				switch txType {
				case constant.TxTypeCoin:
					coinTx, _ := tx.(document.CoinTx)
					coinTx.Height = block.Block.Height
					coinTx.Time = block.Block.Time
					handle(coinTx, mutex, funcChain)
					break
				case stake.TypeTxDeclareCandidacy:
					stakeTxDeclareCandidacy, _ := tx.(document.StakeTxDeclareCandidacy)
					stakeTxDeclareCandidacy.Height = block.Block.Height
					stakeTxDeclareCandidacy.Time = block.Block.Time
					handle(stakeTxDeclareCandidacy, mutex, funcChain)
					break
				case stake.TypeTxEditCandidacy:
					break
				case stake.TypeTxDelegate, stake.TypeTxUnbond:
					stakeTx, _ := tx.(document.StakeTx)
					stakeTx.Height = block.Block.Height
					stakeTx.Time = block.Block.Time
					handle(stakeTx, mutex, funcChain)
					break
				}
			}
		}

		// save block info
		bk := document.Block{
			Height: block.Block.Height,
			Time:   block.Block.Time,
			TxNum:  block.BlockMeta.Header.NumTxs,
		}
		if err := store.Save(bk); err != nil {
			logger.Error.Printf("Save block info failed, err is %v",
				err.Error())
		}
	}
	
	logger.Info.Printf("ThreadNo[%d] finish sync block from %d to %d\n",
		threadNum, start, end)
	
}
