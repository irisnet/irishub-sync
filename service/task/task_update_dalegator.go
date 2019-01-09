package task

import (
	"github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/service/handler"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
)

func MakeUpdateDelegatorTask() Task {
	return NewLockTaskFromEnv(server.CronUpdateDelegator, "save_update_delegator_lock", func() {
		logger.Debug("========================task's trigger [MakeUpdateDelegatorTask] begin===================")
		updateDelegator()
		logger.Debug("========================task's trigger [MakeUpdateDelegatorTask] end===================")
	})
}

func updateDelegator() {
	var delegatorStore document.Delegator
	delegators := delegatorStore.QueryUnbonding()
	if len(delegators) == 0 {
		logger.Info("no delegator is unbonding")
		return
	}

	for _, d := range delegators {
		ubd := handler.BuildUnbondingDelegation(d.Address, d.ValidatorAddr)
		d.UnbondingDelegation = ubd
		store.Update(d)
		logger.Info("update delegator", logger.Any("ubd", ubd))
	}
}
