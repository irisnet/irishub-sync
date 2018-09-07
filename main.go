package main

import (
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/service"
	"github.com/irisnet/irishub-sync/store"
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
	//监听指定信号
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//#########################开启数据库服务##########################
	logger.Info.Println("#########################开启数据库服务##########################")
	store.Start()
	//#########################开启同步服务##########################
	logger.Info.Println("#########################开启同步服务##########################")
	engine.Start()
	//阻塞直至有信号传入
	<-c
}
