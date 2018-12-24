package document

import (
	"encoding/json"
	"testing"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func TestSyncTask_GetMaxBlockHeight(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test get max block height",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := SyncTask{}
			res, err := d.GetMaxBlockHeight()
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("max end_height is %v\n", res)
		})
	}
}

func TestSyncTask_QueryAll(t *testing.T) {
	type args struct {
		status   []string
		taskType string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test query sync task",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := SyncTask{}
			res, err := d.QueryAll(tt.args.status, tt.args.taskType)
			if err != nil {
				t.Fatal(err)
			}
			resBytes, err := json.MarshalIndent(res, "", "\t")
			if err != nil {
				t.Fatalf("marshal json err: %v\n", err)
			}
			t.Log(string(resBytes))
		})
	}
}

func TestSyncTask_GetExecutableTask(t *testing.T) {
	type args struct {
		t int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test get executable task",
			args: args{
				t: 600,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := SyncTask{}
			res, err := d.GetExecutableTask(tt.args.t)
			if err != nil {
				t.Fatal(err)
			}
			resBytes, err := json.MarshalIndent(res, "", "\t")
			if err != nil {
				t.Error(err)
			}
			t.Logf("res is %v\n", string(resBytes))
		})
	}
}

func TestSyncTask_TakeOverTask(t *testing.T) {
	type args struct {
		task     SyncTask
		workerId string
	}
	var (
		syncTaskModel SyncTask
	)

	task1, _ := syncTaskModel.GetTaskById(bson.ObjectIdHex("5c176dc63b6c5c4027b8fb92"))

	tests := []struct {
		name string
		args args
	}{
		{
			name: "test take over task",
			args: args{
				task:     task1,
				workerId: bson.NewObjectIdWithTime(time.Now()).String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := SyncTask{}
			err := d.TakeOverTask(tt.args.task, tt.args.workerId)
			if err != nil {
				if err == mgo.ErrNotFound {
					t.Log("this task has been take over by other goroutine")
				} else {
					t.Fatal(err)
				}
			}
			t.Log("take over task success")
		})
	}
}
