package service

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/service/handler"
	"github.com/irisnet/irishub-sync/service/task"
	"github.com/robfig/cron"
	"time"
)

var (
	engine *SyncEngine
)

func init() {
	engine = &SyncEngine{
		cron:      cron.New(),
		tasks:     []task.Task{},
		initFuncs: []func(){},
	}

	//engine.AddTask(task.MakeSyncBlockTask())
	engine.AddTask(task.MakeCalculateAndSaveValidatorUpTimeTask())
	engine.AddTask(task.MakeCalculateTxGasAndGasPriceTask())
	engine.AddTask(task.MakeSyncProposalStatusTask())
	engine.AddTask(task.MakeValidatorHistoryTask())

	// init delegator for genesis validator
	engine.initFuncs = append(engine.initFuncs, handler.InitDelegator)
}

type SyncEngine struct {
	cron      *cron.Cron  //cron
	tasks     []task.Task // my timer task
	initFuncs []func()    // module init fun
}

func (engine *SyncEngine) AddTask(task task.Task) {
	engine.tasks = append(engine.tasks, task)
	engine.cron.AddFunc(task.GetCron(), task.GetCommand())
}

func (engine *SyncEngine) Start() {
	// init module info
	for _, init := range engine.initFuncs {
		init()
	}
	//watcher := task.NewWatcher()
	//watcher.FastSync()
	go task.StartCreateTask()
	go task.StartExecuteTask()

	// cron task should start after fast sync finished
	fastSyncChan := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			<-ticker.C
			flag, err := task.AssertFastSyncFinished()
			if err != nil {
				logger.Error("assert fast sync finished failed", logger.String("err", err.Error()))
			}
			if flag {
				close(fastSyncChan)
				return
			}
		}
	}()
	<-fastSyncChan
	logger.Info("fast sync finished, now cron task can start")

	engine.cron.Start()
}

func (engine *SyncEngine) Stop() {
	logger.Info("release resource :SyncEngine")
	engine.cron.Stop()
	for _, t := range engine.tasks {
		t.Release()
	}
}

func New() *SyncEngine {
	return engine
}
