package document

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNameSyncConf = "sync_conf"
)

type SyncConf struct {
	BlockNumPerWorker  int64     `bson:"block_num_per_worker"`
	MaxWorkerSleepTime time.Time `bson:"max_worker_sleep_time"`
}

func (d SyncConf) Name() string {
	return CollectionNameSyncConf
}

func (d SyncConf) PkKvPair() map[string]interface{} {
	return bson.M{}
}
