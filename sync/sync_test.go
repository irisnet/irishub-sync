package sync

import (
	"testing"
	"time"
	"fmt"
	"github.com/robfig/cron"
	
	conf "github.com/irisnet/iris-sync-server/conf/server"
	
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	"sync"
	"github.com/irisnet/iris-sync-server/module/logger"
)

func TestStart(t *testing.T) {
	Start()
	for true {
		time.Sleep(time.Minute)
		fmt.Printf("wait\n")
	}
}

func Test_startCron(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	
	spec := conf.SyncCron
	c := cron.New()
	c.AddFunc(spec, func() {
		logger.Info.Println("print word every second")
	})
	go c.Start()
	
	wg.Wait()
}

func Test_watchBlock(t *testing.T) {
	type args struct {
		c rpcClient.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := watchBlock(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("watchBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fastSync(t *testing.T) {
	type args struct {
		c rpcClient.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fastSync(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("fastSync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
