package handler

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
)

func getProposal(cdc *wire.Codec, proposalID int64) (proposal document.Proposal, err error) {
	res, err := helper.Query(gov.KeyProposal(proposalID), "gov", constant.StoreDefaultEndPath)
	if len(res) == 0 || err != nil {
		return proposal, err
	}
	var propo gov.Proposal
	cdc.MustUnmarshalBinary(res, &propo)
	proposal.ProposalId = proposalID
	proposal.Title = propo.GetTitle()
	proposal.Type = propo.GetProposalType().String()
	proposal.Description = propo.GetDescription()
	proposal.Status = propo.GetStatus().String()
	proposal.SubmitBlock = propo.GetSubmitBlock()
	proposal.VotingStartBlock = propo.GetVotingStartBlock()
	proposal.TotalDeposit = types.BuildCoins(propo.GetTotalDeposit())

	votes, err := getVotes(cdc, proposalID)
	if err == nil {
		proposal.Vote = votes
	}
	return
}

func getVotes(cdc *wire.Codec, proposalID int64) (pVotes []document.PVote, err error) {
	res, err := helper.QuerySubspace(cdc, gov.KeyVotesSubspace(proposalID), "gov")
	if len(res) == 0 || err != nil {
		return pVotes, err
	}
	for i := 0; i < len(res); i++ {
		var vote gov.Vote
		cdc.MustUnmarshalBinary(res[i].Value, &vote)
		v := document.PVote{
			Voter:  vote.Voter.String(),
			Option: vote.Option.String(),
		}
		pVotes = append(pVotes, v)
	}
	return
}

func handleProposal(docTx document.CommonTx) {
	switch docTx.Type {
	case constant.TxTypeSubmitProposal:
		if proposal, err := getProposal(codec.Cdc, docTx.ProposalId); err == nil {
			proposal.SubmitTime = docTx.Time
			store.Save(proposal)
		}
	case constant.TxTypeDeposit, constant.TxTypeVote:
		if proposal, err := document.QueryProposal(docTx.ProposalId); err == nil {
			propo, _ := getProposal(codec.Cdc, docTx.ProposalId)
			propo.SubmitTime = proposal.SubmitTime
			store.SaveOrUpdate(propo)
		}
	}
}
