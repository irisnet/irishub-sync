package handler

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
)

func handleProposal(docTx document.CommonTx) {
	switch docTx.Type {
	case constant.TxTypeSubmitProposal:
		if proposal, err := helper.GetProposal(docTx.ProposalId); err == nil {
			if isContainVotingPeriodStartTag(docTx) {
				proposal.VotingPeriodStartHeight = docTx.Height
			}
			store.SaveOrUpdate(proposal)
		}
	case constant.TxTypeDeposit:
		if proposal, err := document.QueryProposal(docTx.ProposalId); err == nil {
			propo, _ := helper.GetProposal(docTx.ProposalId)
			if isContainVotingPeriodStartTag(docTx) {
				proposal.VotingPeriodStartHeight = docTx.Height
			}
			proposal.TotalDeposit = propo.TotalDeposit
			proposal.Status = propo.Status
			proposal.VotingStartTime = propo.VotingStartTime
			proposal.VotingEndTime = propo.VotingEndTime
			store.SaveOrUpdate(proposal)
		}
	case constant.TxTypeVote:
		//失败的投票不计入统计
		if docTx.Status == document.TxStatusFail {
			return
		}
		if proposal, err := document.QueryProposal(docTx.ProposalId); err == nil {
			voteMsg := docTx.Msg.(types.Vote)
			vote := document.PVote{
				Voter:  voteMsg.Voter,
				Option: voteMsg.Option,
				Time:   docTx.Time,
			}
			var i int
			var hasVote = false
			for i = range proposal.Votes {
				if proposal.Votes[i].Voter == vote.Voter {
					hasVote = true
					break
				}
			}
			if hasVote {
				proposal.Votes[i] = vote
			} else {
				proposal.Votes = append(proposal.Votes, vote)
			}
			store.SaveOrUpdate(proposal)
		}
	}
}

func isContainVotingPeriodStartTag(docTx document.CommonTx) bool {
	tags := docTx.Tags
	if len(tags) > 0 {
		for k, _ := range tags {
			if k == constant.TxTagVotingPeriodStart {
				return true
			}
		}
	}

	return false
}
