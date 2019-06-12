package task

import (
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
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
				propo.SubmitTime = proposal.SubmitTime
				propo.Votes = proposal.Votes
				propo.VotingPeriodStartHeight = proposal.VotingPeriodStartHeight
				store.SaveOrUpdate(propo)
			}
		}
	}
}

func MakeSyncProposalStatusTask() Task {
	return NewLockTaskFromEnv(conf.SyncProposalStatus, func() {
		logger.Debug("========================task's trigger [SyncProposalStatus] begin===================")
		syncProposalStatus()
		logger.Debug("========================task's trigger [SyncProposalStatus] end===================")
	})
}
