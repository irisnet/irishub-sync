package document

import (
	"encoding/json"
	"testing"
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
