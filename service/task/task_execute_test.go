package task

import (
	"sync"
	"testing"

	"encoding/json"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/helper"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func Test_executeTask(t *testing.T) {
	type args struct {
		blockNumPerWorkerHandle int64
		maxWorkerSleepTime      int64
		chanLimit               chan bool
	}

	limitChan := make(chan bool, 2)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "test execute task",
			args: args{
				blockNumPerWorkerHandle: 100,
				maxWorkerSleepTime:      10 * 60,
				chanLimit:               limitChan,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)

			tt.args.chanLimit <- true
			go executeTask(tt.args.blockNumPerWorkerHandle, tt.args.maxWorkerSleepTime, tt.args.chanLimit)

			wg.Wait()
		})
	}
}

func Test_assertTaskValid(t *testing.T) {
	var (
		syncTaskModel document.SyncTask
	)

	catchUpTask, _ := syncTaskModel.GetTaskById(bson.ObjectIdHex("5c176b243b6c5c3ff62deaea"))
	followTask, _ := syncTaskModel.GetTaskById(bson.ObjectIdHex("5c176dc63b6c5c4027b8fb92"))

	type args struct {
		task                    document.SyncTask
		blockNumPerWorkerHandle int64
		blockChainLatestHeight  int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test assert catch up task valid",
			args: args{
				task:                    catchUpTask,
				blockNumPerWorkerHandle: 100,
				blockChainLatestHeight:  100,
			},
		},
		{
			name: "test assert follow task valid",
			args: args{
				task:                    followTask,
				blockNumPerWorkerHandle: 200,
				blockChainLatestHeight:  800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assertTaskValid(tt.args.task, tt.args.blockNumPerWorkerHandle, tt.args.blockChainLatestHeight)
			t.Log(got)
		})
	}
}

func Test_parseBlock(t *testing.T) {
	client := helper.GetClient()

	defer func() {
		logger.Info("release client")
		client.Release()
	}()

	type args struct {
		b      int64
		client *helper.Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test parse block",
			args: args{
				client: client,
				b:      107061,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parseBlock(tt.args.b, tt.args.client)
			if err != nil {
				t.Fatal(err)
			}
			resBytes, err := json.MarshalIndent(res, "", "\t")
			if err != nil {
				t.Error(err)
			}
			t.Log(string(resBytes))
		})
	}
}

func Test_assertTaskWorkerUnchanged(t *testing.T) {
	type args struct {
		taskId   bson.ObjectId
		workerId string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test assert test worker unchanged",
			args: args{
				taskId:   bson.ObjectIdHex("5b176dc63b6c5c4027b8fb92"),
				workerId: "ObjectIdHex(\"5c19cec10000000000000000\")",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := assertTaskWorkerUnchanged(tt.args.taskId, tt.args.workerId)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(res)
		})
	}
}

func Test_saveDocs(t *testing.T) {
	var (
		syncTaskModel document.SyncTask
	)

	block := document.Block{
		Height: 480,
		Hash:   bson.NewObjectId().Hex(),
	}
	task, _ := syncTaskModel.GetTaskById(bson.ObjectIdHex("5c19e9e03b6c5ca9a96fdf62"))
	task.CurrentHeight = block.Height
	task.LastUpdateTime = time.Now().Unix()

	type args struct {
		blockDoc document.Block
		taskDoc  document.SyncTask
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test save docs",
			args: args{
				blockDoc: block,
				taskDoc:  task,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := saveDocs(tt.args.blockDoc, tt.args.taskDoc); err != nil {
				t.Fatal(err)
			}
			t.Log("save success")
		})
	}
}
