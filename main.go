package main

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/service"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/helper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal)
	engine := service.New()

	defer func() {
		logger.Info("#########################System Exit##########################")
		engine.Stop()
		helper.ClosePool()
		store.Stop()
		logger.Sync()
		if err := recover(); err != nil {
			logger.Error(err.(string))
			os.Exit(1)
		}
	}()
	//监听指定信号
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//#########################开启数据库服务##########################
	logger.Info("#########################开启数据库服务##########################")
	store.Start()
	//#########################开启同步服务##########################
	logger.Info("#########################开启同步服务##########################")
	engine.Start()
	//阻塞直至有信号传入
	<-c
}
