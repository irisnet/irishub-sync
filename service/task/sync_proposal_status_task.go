package task

import (
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
)

func syncProposalStatus() {
	var status = []string{constant.StatusDepositPeriod, constant.StatusVotingPeriod}
	if proposals, err := document.QueryByStatus(status); err == nil {
		for _, proposal := range proposals {
			propo, err := helper.GetProposal(proposal.ProposalId)
			if err != nil {
				store.Delete(proposal)
				return
			}
			if propo.Status != proposal.Status {
				proposal.Status = propo.Status
				store.SaveOrUpdate(proposal)
			}
		}
	}
}

func MakeSyncProposalStatusTask() Task {
	return NewLockTaskFromEnv(conf.SyncProposalStatus, "sync_proposal_status_lock", func() {
		logger.Info.Printf("========================task's trigger [%s] begin===================", "SyncProposalStatus")
		syncProposalStatus()
		logger.Info.Printf("========================task's trigger [%s] end===================", "SyncProposalStatus")
	})
}
