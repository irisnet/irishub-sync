package main

import (
	"github.com/irisnet/iris-sync-server/sync"
	"time"
	"github.com/irisnet/iris-sync-server/module/logger"
)

func main() {
	sync.Start()
	for true {
		time.Sleep(time.Minute)
		logger.Info.Println("wait")
	}
}
