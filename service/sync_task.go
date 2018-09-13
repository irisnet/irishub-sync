package service

import (
	"github.com/irisnet/irishub-sync/util/helper"
	"time"
)

type Command = func()
type Task struct {
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
	return Task{
		Spec:     spec,
		LockKey:  lockKey,
		cmd:      cmd,
		lock:     lock,
		withLock: withLock,
	}
}

func (task Task) Stop() {
	if task.withLock {
		task.lock.Destroy()
	}
}

func (task Task) before() {
	if task.withLock {
		task.lock.Destroy()
	}
}

func (task Task) after() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	if task.withLock {
		task.lock.UnLock()
		task.lock.Destroy()
	}
}

func (task Task) GetCmd() Command {
	lockCmd := func() {
		task.before()
		task.cmd()
		task.after()
	}
	return lockCmd
}
