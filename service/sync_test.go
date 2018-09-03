package service

import (
	"github.com/robfig/cron"
	"testing"

	conf "github.com/irisnet/irishub-sync/conf/server"

	"github.com/irisnet/irishub-sync/module/logger"
	"sync"
)

func TestStart(t *testing.T) {
	var (
		limitChan    chan int
		unBufferChan chan int
	)
	limitChan = make(chan int, 3)
	unBufferChan = make(chan int)
	goroutineNum := 5
	activeGoroutineNum := goroutineNum
	for i := 1; i <= goroutineNum; i++ {
		limitChan <- i
		go func(goroutineNum int, ch chan int) {
			logger.Info.Println("release limitChan")
			<-limitChan
			defer func() {
				logger.Info.Printf("%v goroutine send data to channel\n",
					goroutineNum)
				ch <- goroutineNum
			}()

		}(i, nil)
	}

	for {
		select {
		case <-unBufferChan:
			activeGoroutineNum = activeGoroutineNum - 1
			logger.Info.Printf("active goroutine num is %v", activeGoroutineNum)
			if activeGoroutineNum == 0 {
				logger.Info.Println("All goroutine complete")
				break
			}
		}
	}

}

func Test_startCron(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	c := cron.New()
	c.AddFunc(conf.CronCalculateUpTime, func() {
		logger.Info.Println("every one minute execute code")
	})
	c.AddFunc(conf.CronCalculateTxGas, func() {
		logger.Info.Println("every five minute execute code")
	})
	go c.Start()

	wg.Wait()
}
