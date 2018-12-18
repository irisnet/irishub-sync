package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNameSyncConf = "sync_conf"
)

type SyncConf struct {
	ChainId                 string `bson:"chain_id"`
	BlockNumPerWorkerHandle int64  `bson:"block_num_per_worker_handle"`
	MaxWorkerSleepTime      int64  `bson:"max_worker_sleep_time"`
}

func (d SyncConf) Name() string {
	return CollectionNameSyncConf
}

func (d SyncConf) PkKvPair() map[string]interface{} {
	return bson.M{"chain_id": d.ChainId}
}

func (d SyncConf) GetConf(chainId string) (SyncConf, error) {
	var syncConf SyncConf

	q := bson.M{
		"chain_id": chainId,
	}
	fn := func(c *mgo.Collection) error {
		return c.Find(q).One(&syncConf)
	}

	err := store.ExecCollection(d.Name(), fn)

	if err != nil {
		return syncConf, err
	}

	return syncConf, nil
}
