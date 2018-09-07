package main

import (
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/service"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/helper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal)
	engine := service.GetSyncEngine()

	defer func() {
		logger.Info.Println("#########################System Exit##########################")
		engine.Stop()
		helper.ClosePool()
		store.Stop()
		if err := recover(); err != nil {
			logger.Error.Println(err)
			os.Exit(1)
		}
	}()
	//监听指定信号 ctrl+c kill
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//#########################开启数据库服务##########################
	logger.Info.Println("#########################开启数据库服务##########################")
	startDb()
	//#########################开启同步服务##########################
	logger.Info.Println("#########################开启同步服务##########################")
	engine.Start()
	//阻塞直至有信号传入
	<-c
}

func startDb() {
	store.Start()
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
}
