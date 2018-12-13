package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNameSyncConf = "sync_conf"
)

type SyncConf struct {
	BlockNumPerWorkerHandle int64     `bson:"block_num_per_worker_handle"`
	MaxWorkerSleepTime      time.Time `bson:"max_worker_sleep_time"`
}

func (d SyncConf) Name() string {
	return CollectionNameSyncConf
}

func (d SyncConf) PkKvPair() map[string]interface{} {
	return bson.M{}
}

func (d SyncConf) GetConf() (SyncConf, error) {
	var syncConf SyncConf

	q := bson.M{}
	fn := func(c *mgo.Collection) error {
		return c.FindId(q).One(&syncConf)
	}

	err := store.ExecCollection(d.Name(), fn)

	if err != nil {
		return syncConf, err
	}

	return syncConf, nil
}
