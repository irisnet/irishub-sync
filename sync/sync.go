package sync

import (
	conf "github.com/irisnet/iris-sync-server/conf/server"
	"github.com/irisnet/iris-sync-server/model/store"
	"github.com/irisnet/iris-sync-server/module/logger"
	"github.com/irisnet/iris-sync-server/util/helper"
	"github.com/irisnet/iris-sync-server/util/constant"

	"github.com/robfig/cron"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	"github.com/irisnet/iris-sync-server/model/store/document"
	"strings"
	"encoding/hex"
)

var (
	// how many block each goroutine need to sync when do fast sync
	syncBlockNumFastSync = int64(conf.SyncBlockNumFastSync)

	// limit max goroutine
	limitChan = make(chan int64, conf.SyncMaxGoroutine)
)

// start sync server
func Start() {
	InitServer()
	c := helper.GetClient().Client
	if err := fastSync(c); err != nil {
		logger.Error.Fatalf("sync block failed,%v\n", err)
	}
	startCron(c)
	//go watch(c) 监控的方式在启动同步过程中容易丢失区块
}

func InitServer() {
	store.Init()

	chainId := conf.ChainId
	syncTask, err := document.QuerySyncTask()

	if err != nil {
		if chainId == "" {
			logger.Error.Fatalln("sync process start failed,chainId is empty\n")
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
	b, _ := document.QuerySyncTask()
	status, _ := c.Status()
	latestBlockHeight := status.LatestBlockHeight

	funcChain := []func(tx store.Docs){
		saveTx, saveOrUpdateAccount, updateAccountBalance,
	}

	ch := make(chan int64)
	limitChan <- 1

	go syncBlock(b.Height+1, latestBlockHeight, funcChain, ch, 0)

	select {
	case <-ch:
		//更新同步记录
		block, _ := c.Block(&latestBlockHeight)
		b.Height = block.Block.Height
		b.Time = block.Block.Time
		return store.Update(b)
	}
}

// fast sync data from blockChain
func fastSync(c rpcClient.Client) error {
	b, _ := document.QuerySyncTask()
	status, _ := c.Status()
	latestBlockHeight := status.LatestBlockHeight
	funcChain := []func(tx store.Docs){
		saveTx, saveOrUpdateAccount,
	}

	ch := make(chan int64)
	activeThreadNum := int64(0)

	goRoutineNum := (latestBlockHeight - b.Height) / syncBlockNumFastSync

	if goRoutineNum == 0 {
		goRoutineNum = 10
		syncBlockNumFastSync = 100
	}

	for i := int64(1); i <= goRoutineNum; i++ {
		activeThreadNum++
		limitChan <- i
		var (
			start = b.Height + (i-1)*syncBlockNumFastSync + 1
			end   = b.Height + i*syncBlockNumFastSync
		)
		if i == goRoutineNum {
			end = latestBlockHeight
		}
		go syncBlock(start, end, funcChain, ch, i)
	}

	//threadNum := (latestBlockHeight - b.Height) / maxBatchNum
	//// 单线程处理
	//if threadNum == 0 {
	//	go syncBlock(b.Height, latestBlockHeight, funcChain, ch, 0)
	//} else {
	//	// 开启多线程处理
	//	for i := int64(1); i <= threadNum; i++ {
	//		activeThreadNum ++
	//		var start = b.Height + (i-1)*maxBatchNum + 1
	//		var end = b.Height + i*maxBatchNum
	//		if i == threadNum {
	//			end = latestBlockHeight
	//		}
	//
	//		go syncBlock(start, end, funcChain, ch, i)
	//	}
	//
	//}

	for {
		select {
		case threadNo := <-ch:
			logger.Info.Printf("threadNo[%d] is over\n", threadNo)
			activeThreadNum = activeThreadNum - 1
			if activeThreadNum == 0 {
				goto end
			}
		}
	}

end:
	{
		//更新同步记录
		block, _ := c.Block(&latestBlockHeight)
		b.Height = block.Block.Height
		b.Time = block.Block.Time
		store.Update(b)

		//同步账户余额
		accounts := document.QueryAll()
		for _, account := range accounts {
			updateAccountBalance(account)
		}
		logger.Info.Println("update account balance over")

		return nil
	}
}

func syncBlock(start int64, end int64, funcChain []func(tx store.Docs), ch chan int64, threadNum int64) {
	logger.Info.Printf("threadNo[%d] begin sync block from %d to %d\n", threadNum, start, end)
	client := helper.GetClient()
	// release client
	defer client.Release()
	// release buffer chain
	defer func() {
		<-limitChan
	}()

	for j := start; j <= end; j++ {
		logger.Info.Printf("===========threadNo[%d] sync block,height:%d===========\n", threadNum, j)

		// TODO 使用client.Client.BlockChainInfo
		block, err := client.Client.Block(&j)
		if err != nil {
			// try again
			client2 := helper.GetClient()
			block, err = client2.Client.Block(&j)
			if err != nil {
				logger.Error.Fatalf("invalid block height %d\n", j)
			}
		}
		if block.BlockMeta.Header.NumTxs > 0 {
			txs := block.Block.Data.Txs
			for _, txByte := range txs {
				txType, tx := helper.ParseTx(txByte)
				txHash := strings.ToUpper(hex.EncodeToString(txByte.Hash()))
				logger.Info.Printf("===========threadNo[%d] find tx,txType=%s;txHash=%s\n", threadNum, txType, txHash)

				switch txType {
				case constant.TxTypeCoin:
					coinTx, _ := tx.(document.CoinTx)
					coinTx.Height = block.Block.Height
					coinTx.Time = block.Block.Time
					handle(coinTx, funcChain)
					break
				case constant.TxTypeStake:
					stakeTx, _ := tx.(document.StakeTx)
					stakeTx.Height = block.Block.Height
					stakeTx.Time = block.Block.Time
					handle(stakeTx, funcChain)
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
		store.SaveOrUpdate(bk)

	}
	ch <- threadNum
}
