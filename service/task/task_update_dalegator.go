package task

import (
	"github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/service/handler"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
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

	var ops []txn.Op
	for _, d := range delegators {
		unb := handler.BuildUnbondingDelegation(d.Address, d.ValidatorAddr)
		delegation := handler.BuildDelegation(d.Address, d.ValidatorAddr)
		updateOp := txn.Op{
			C:      document.CollectionNmStakeRoleDelegator,
			Id:     d.ID,
			Assert: txn.DocExists,
			Update: bson.M{
				"$set": bson.M{
					"unbonding_delegation": unb,
					"shares":               delegation.Shares,
					"original_shares":      delegation.OriginalShares,
					"height":               delegation.Height,
				},
			},
		}
		ops = append(ops, updateOp, updateOp)
		logger.Info("update delegator", logger.Any("updateOp", updateOp))
	}
	if len(ops) > 0 {
		err := store.Txn(ops)
		if err != nil {
			logger.Info("update delegator failed")
			return
		}
	}

}
