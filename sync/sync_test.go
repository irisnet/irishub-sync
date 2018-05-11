package sync

import (
	"testing"

	"github.com/irisnet/iris-sync-server/model/store"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	"fmt"
	"time"
)

func TestStart(t *testing.T) {
	Start()
	for true {
		time.Sleep(time.Minute)
		fmt.Printf("wait\n")
	}
}

func Test_startCron(t *testing.T) {
	type args struct {
		client rpcClient.Client
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startCron(tt.args.client)
		})
	}
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

func Test_syncBlock(t *testing.T) {
	type args struct {
		start     int64
		end       int64
		funcChain []func(tx store.Docs)
		ch        chan int64
		threadNum int64
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			syncBlock(tt.args.start, tt.args.end, tt.args.funcChain, tt.args.ch, tt.args.threadNum)
		})
	}
}
