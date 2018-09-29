package task

import (
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/util/helper"
	"time"
)

type Command = func()
type Task interface {
	GetCommand() Command
	GetCron() string
	Release()
}

type LockTask struct {
	Spec     string
	LockKey  string
	cmd      Command
	lock     *helper.DLock
	withLock bool
}

func NewTask(spec, lockKey string, cmd Command, withLock bool) Task {
	var lock *helper.DLock
	if withLock {
		if len(lockKey) == 0 {
			panic("lockKey can not be empty")
		}
		lock = helper.NewLock(lockKey, 500*time.Millisecond)
	}
	return &LockTask{
		Spec:     spec,
		LockKey:  lockKey,
		cmd:      cmd,
		lock:     lock,
		withLock: withLock,
	}
}

func NewLockTaskFromEnv(spec, lockKey string, cmd Command) Task {
	return NewTask(spec, lockKey, cmd, conf.SyncWithDLock)
}

func (task *LockTask) before() {
	if task.withLock {
		task.lock.Lock()
	}
}

func (task *LockTask) after() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	if task.withLock {
		task.lock.UnLock()
	}
}

func (task *LockTask) GetCommand() Command {
	lockCmd := func() {
		task.before()
		task.cmd()
		task.after()
	}
	return lockCmd
}

func (task *LockTask) GetCron() string {
	return task.Spec
}

func (task *LockTask) Release() {
	if task.withLock {
		task.lock.Destroy()
	}
}
