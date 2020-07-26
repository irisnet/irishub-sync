package handler

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
	"strconv"
	. "github.com/irisnet/irishub-sync/msg"
	"github.com/irisnet/irishub-sync/types"
	"encoding/json"
)

func handleProposal(docTx *document.CommonTx) {
	switch docTx.Type {
	case constant.TxTypeSubmitProposal:
		if proposal, err := helper.GetProposal(docTx.ProposalId); err == nil {
			if isContainVotingPeriodStartEvent(docTx) {
				proposal.VotingPeriodStartHeight = docTx.Height
			}
			proposal.Type = getProposalTypeFromEvents(docTx.Events)

			store.SaveOrUpdate(proposal)
		}
	case constant.TxTypeDeposit:
		if proposal, err := document.QueryProposal(docTx.ProposalId); err == nil {
			propo, _ := helper.GetProposal(docTx.ProposalId)
			if isContainVotingPeriodStartEvent(docTx) {
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
			msgVote := DocTxMsgVote{}
			data, _ := json.Marshal(docTx.Msgs[0].Msg)
			json.Unmarshal(data, &msgVote)
			vote := document.PVote{
				Voter:  msgVote.Voter,
				Option: msgVote.Option,
				TxHash: docTx.TxHash,
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

func isContainVotingPeriodStartEvent(docTx *document.CommonTx) (bool) {
	events := docTx.Events
	if len(events) > 0 {
		for _, one := range events {
			if one.Type != types.EventTypeProposalDeposit {
				continue
			}
			for k, _ := range one.Attributes {
				if k == types.EventGovVotingPeriodStart {
					return true
				}
			}

		}
	}

	return false
}

func IsContainVotingEndEvent(events []document.Event) (uint64, bool) {
	//events := blockresult.Events
	if len(events) > 0 {
		for _, event := range events {
			if val, ok := event.Attributes[types.EventGovProposalID]; ok {
				proposalid, _ := strconv.ParseUint(val, 10, 64)
				return proposalid, true
			}
		}
	}
	return 0, false
}

func getProposalTypeFromEvents(result []document.Event) (string) {
	//query proposal type
	for _, val := range result {
		if val.Type != types.EventTypeSubmitProposal {
			continue
		}
		for key, val := range val.Attributes {
			if key == types.EventGovProposalType {
				return val
			}
		}
	}

	return ""
}
