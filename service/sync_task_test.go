package service

import (
	"github.com/irisnet/irishub-sync/store"
	"os"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	store.Start()

	code := m.Run()

	os.Exit(code)
}

func Test_createTask(t *testing.T) {
	type args struct {
		blockNumPerWorker int64
		chanLimit         chan bool
	}

	chanLimit := make(chan bool, 2)
	var wg sync.WaitGroup

	tests := []struct {
		name string
		args args
	}{
		{
			name: "test create task",
			args: args{
				blockNumPerWorker: 100,
				chanLimit:         chanLimit,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg.Add(1)
			chanLimit <- true
			go createTask(tt.args.blockNumPerWorker, tt.args.chanLimit)
			wg.Wait()
		})
	}
}
