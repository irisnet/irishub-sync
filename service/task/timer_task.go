package task

type Command = func()
type Task interface {
	GetCommand() Command
	GetCron() string
}

type LockTask struct {
	Spec     string
	cmd      Command
	withLock bool
}

func NewTask(spec string, cmd Command) Task {
	return &LockTask{
		Spec: spec,
		cmd:  cmd,
	}
}

func NewLockTaskFromEnv(spec string, cmd Command) Task {
	return NewTask(spec, cmd)
}

func (task *LockTask) GetCommand() Command {
	lockCmd := func() {
		task.cmd()
	}
	return lockCmd
}

func (task *LockTask) GetCron() string {
	return task.Spec
}
