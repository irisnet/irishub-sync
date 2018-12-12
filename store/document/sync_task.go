package document

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNameSyncTask = "sync_task"
)

type WorkerLog struct {
	WorkerId  int64     `bson:"worker_id"`  // worker id
	BeginTime time.Time `bson:"begin_time"` // time which worker begin handle this task
}

type SyncTask struct {
	StartHeight    int64       `bson:"start_height"`   // task start height
	EndHeight      int64       `bson:"end_height"`     // task end height
	CurrentHeight  int64       `bson:"current_height"` // task current height
	Status         string      `bson:"status"`         // task status
	WorkerId       int64       `bson:"worker_id"`      // worker id
	WorkerLogs     []WorkerLog `bson:"worker_logs"`    // worker logs
	LastUpdateTime time.Time   `bson:"last_update_time"`
}

func (d SyncTask) Name() string {
	return CollectionNameSyncTask
}

func (d SyncTask) PkKvPair() map[string]interface{} {
	return bson.M{"start_height": d.CurrentHeight, "end_height": d.EndHeight}
}
