package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmCronTask = "cron_task"

	CronTaskStatusReady      = "ready"
	CronTaskStatusProcessing = "processing"
)

type CronTask struct {
	Key       string    `bson:"key"`
	Cron      string    `bson:"cron"`
	Worker    string    `bson:"worker"`
	Status    string    `bson:"status"`
	BeginTime time.Time `bson:"begin_time"`
	EndTime   time.Time `bson:"end_time"`
}

func (d CronTask) Name() string {
	return CollectionNmCronTask
}

func (d CronTask) PkKvPair() map[string]interface{} {
	return bson.M{"key": d.Key}
}

func LockCronTask(key, worker string) bool {

	update := func(c *mgo.Collection) error {
		selector := bson.M{
			"key":    key,
			"status": CronTaskStatusReady,
		}
		updator := bson.M{
			"$set": bson.M{
				"status":     CronTaskStatusProcessing,
				"worker":     worker,
				"begin_time": time.Now(),
			},
		}
		return c.Update(selector, updator)
	}

	if err := store.ExecCollection(CollectionNmCronTask, update); err != nil {
		return false
	}
	return true
}

func UnLockCronTask(key, worker string) bool {

	update := func(c *mgo.Collection) error {
		selector := bson.M{
			"key":    key,
			"worker": worker,
			"status": CronTaskStatusProcessing,
		}
		updator := bson.M{
			"$set": bson.M{
				"worker":   "",
				"status":   CronTaskStatusReady,
				"end_time": time.Now(),
			},
		}
		return c.Update(selector, updator)
	}

	if err := store.ExecCollection(CollectionNmCronTask, update); err != nil {
		return false
	}
	return true
}
