package main

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/monitor"
	"github.com/irisnet/irishub-sync/service"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/helper"
	"github.com/irisnet/irishub-sync/cron"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal)
	engine := service.New()

	defer func() {
		logger.Info("Irishub Sync Service Exit...")
		engine.Stop()
		helper.ClosePool()
		store.Stop()
		logger.Sync()
		if err := recover(); err != nil {
			logger.Error(err.(string))
			os.Exit(1)
		}
	}()
	//monitor system signal
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// start monitor
	go monitor.NewMonitor().Start()
	//start databases service
	logger.Info("Databases Service Start...")
	store.Start()
	//start sync task service
	logger.Info("Irishub Sync Service Start...")
	go new(cron.CronService).StartCronService()
	engine.Start()
	//paused until the signal have received
	<-c
}
