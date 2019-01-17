package task

import (
	"fmt"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"gopkg.in/mgo.v2/bson"
	"os"
)

type Command = func()

type CronTask struct {
	cron string
	cmd  Command
	key  string
}

func NewLockTask(key, cron string, cmd Command) CronTask {
	return CronTask{
		cron: cron,
		cmd:  cmd,
		key:  key,
	}
}

func (task *CronTask) GetCmd() Command {
	lockCmd := func() {
		genWorkerId := func() string {
			hostname, _ := os.Hostname()
			return fmt.Sprintf("%v@%v", hostname, bson.NewObjectId().Hex())
		}
		worker := genWorkerId()
		ok := document.LockCronTask(task.GetKey(), worker)
		if ok {
			task.cmd()
			ok = document.UnLockCronTask(task.GetKey(), worker)
			if !ok {
				logger.Error("cron task unlock failed", logger.String("key", task.key), logger.String("worker", worker))
			}
		} else {
			logger.Info("cron task be token ")
		}
	}
	return lockCmd
}

func (task *CronTask) GetCron() string {
	return task.cron
}

func (task *CronTask) GetKey() string {
	return task.key
}

func (task *CronTask) Init() {
	t := document.CronTask{
		Key:    task.key,
		Cron:   task.cron,
		Status: document.CronTaskStatusReady,
	}
	store.Save(t)
}
