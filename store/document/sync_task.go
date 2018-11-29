package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmSyncTask = "sync_task"

	SyncTask_Field_ChainID = "chain_id"
	SyncTask_Field_Height  = "height"
	SyncTask_Field_Time    = "time"
	SyncTask_Field_Syncing = "syncing"
)

type SyncTask struct {
	ChainID string    `bson:"chain_id"`
	Height  int64     `bson:"height"`
	Time    time.Time `bson:"time"`
	Syncing bool      `bson:"syncing"`
}

func (c SyncTask) Name() string {
	return CollectionNmSyncTask
}

func (c SyncTask) PkKvPair() map[string]interface{} {
	return bson.M{SyncTask_Field_ChainID: c.ChainID}
}

func QuerySyncTask() (SyncTask, error) {
	result := SyncTask{}

	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{}).One(&result)
		return err
	}

	err := store.ExecCollection(CollectionNmSyncTask, query)
	return result, err
}
