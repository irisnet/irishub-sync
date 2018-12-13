package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNameSyncTask = "sync_task"

	// value of status
	SyncTaskStatusUnHandled = "unhandled"
	SyncTaskStatusUnderway  = "underway"
	SyncTaskStatusCompleted = "completed"

	// taskType
	SyncTaskTypeCatchUp = "catch_up"
	SyncTaskTypeFollow  = "follow"
)

type WorkerLog struct {
	WorkerId  int64     `bson:"worker_id"`  // worker id
	BeginTime time.Time `bson:"begin_time"` // time which worker begin handle this task
}

type SyncTask struct {
	ID             bson.ObjectId `bson:"_id"`
	StartHeight    int64         `bson:"start_height"`   // task start height
	EndHeight      int64         `bson:"end_height"`     // task end height
	CurrentHeight  int64         `bson:"current_height"` // task current height
	Status         string        `bson:"status"`         // task status
	WorkerId       int64         `bson:"worker_id"`      // worker id
	WorkerLogs     []WorkerLog   `bson:"worker_logs"`    // worker logs
	LastUpdateTime time.Time     `bson:"last_update_time"`
}

func (d SyncTask) Name() string {
	return CollectionNameSyncTask
}

func (d SyncTask) PkKvPair() map[string]interface{} {
	return bson.M{"start_height": d.CurrentHeight, "end_height": d.EndHeight}
}

// get max block height in sync task
func (d SyncTask) GetMaxBlockHeight() (int64, error) {
	type maxHeightRes struct {
		MaxHeight int64 `bson:"max"`
	}
	var res []maxHeightRes

	q := []bson.M{
		{
			"$group": bson.M{
				"_id": nil,
				"max": bson.M{"$sum": "$end_height"},
			},
		},
	}

	getMaxBlockHeightFn := func(c *mgo.Collection) error {
		return c.Pipe(q).All(&res)
	}
	err := store.ExecCollection(d.Name(), getMaxBlockHeightFn)

	if err != nil {
		return 0, err
	}
	if len(res) > 0 {
		return res[0].MaxHeight, nil
	}

	return 0, nil
}

// query record by status
func (d SyncTask) QueryAll(status []string, taskType string) ([]SyncTask, error) {
	var syncTasks []SyncTask
	q := bson.M{}

	if len(status) > 0 {
		q["status"] = bson.M{
			"$in": status,
		}
	}

	switch taskType {
	case SyncTaskTypeCatchUp:
		q["end_height"] = bson.M{
			"$ne": 0,
		}
		break
	case SyncTaskTypeFollow:
		q["end_height"] = bson.M{
			"$eq": 0,
		}
		break
	}

	fn := func(c *mgo.Collection) error {
		return c.Find(q).All(&syncTasks)
	}

	err := store.ExecCollection(d.Name(), fn)

	if err != nil {
		return syncTasks, err
	}

	return syncTasks, nil
}
