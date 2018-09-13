package helper

import (
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

func GetProposal(proposalID int64) (proposal document.Proposal, err error) {
	res, err := Query(gov.KeyProposal(proposalID), "gov", constant.StoreDefaultEndPath)
	if len(res) == 0 || err != nil {
		return proposal, err
	}
	var propo gov.Proposal
	codec.Cdc.MustUnmarshalBinary(res, &propo)
	proposal.ProposalId = proposalID
	proposal.Title = propo.GetTitle()
	proposal.Type = propo.GetProposalType().String()
	proposal.Description = propo.GetDescription()
	proposal.Status = propo.GetStatus().String()
	proposal.SubmitBlock = propo.GetSubmitBlock()
	proposal.VotingStartBlock = propo.GetVotingStartBlock()
	proposal.TotalDeposit = types.BuildCoins(propo.GetTotalDeposit())
	proposal.Votes = []document.PVote{}
	return
}

func GetVotes(proposalID int64) (pVotes []document.PVote, err error) {
	res, err := QuerySubspace(codec.Cdc, gov.KeyVotesSubspace(proposalID), "gov")
	if len(res) == 0 || err != nil {
		return pVotes, err
	}
	for i := 0; i < len(res); i++ {
		var vote gov.Vote
		codec.Cdc.MustUnmarshalBinary(res[i].Value, &vote)
		v := document.PVote{
			Voter:  vote.Voter.String(),
			Option: vote.Option.String(),
		}
		pVotes = append(pVotes, v)
	}
	return
}
