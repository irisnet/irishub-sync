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
			proposal.SubmitTime = docTx.Time
			store.Save(proposal)
		}
	case constant.TxTypeDeposit:
		if proposal, err := document.QueryProposal(docTx.ProposalId); err == nil {
			propo, _ := helper.GetProposal(docTx.ProposalId)
			proposal.TotalDeposit = propo.TotalDeposit
			proposal.Status = propo.Status
			store.SaveOrUpdate(proposal)
		}
	case constant.TxTypeVote:
		if proposal, err := document.QueryProposal(docTx.ProposalId); err == nil {
			proposal.VotingStartBlock = docTx.Height
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
