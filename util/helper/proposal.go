package helper

import (
	"errors"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

func GetProposal(proposalID uint64) (proposal document.Proposal, err error) {
	cdc := types.GetCodec()

	res, err := Query(types.KeyProposal(proposalID), "gov", constant.StoreDefaultEndPath)
	if len(res) == 0 || err != nil {
		return proposal, errors.New("no data")
	}
	var propo types.Proposal
	cdc.UnmarshalBinaryLengthPrefixed(res, &propo) //TODO
	proposal.ProposalId = proposalID
	proposal.Title = propo.GetTitle()
	//proposal.Type = propo.GetProposalType().String()
	//proposal.Description = propo.GetDescription()
	proposal.Status = propo.Status.String()

	proposal.SubmitTime = propo.SubmitTime
	proposal.VotingStartTime = propo.VotingStartTime
	proposal.VotingEndTime = propo.VotingEndTime
	proposal.DepositEndTime = propo.DepositEndTime
	proposal.TotalDeposit = types.ParseCoins(propo.TotalDeposit.String())
	proposal.Votes = []document.PVote{}

	tallyResult := propo.FinalTallyResult
	proposal.TallyResult = document.PTallyResult{
		Yes:               tallyResult.Yes.String(),
		Abstain:           tallyResult.Abstain.String(),
		No:                tallyResult.No.String(),
		NoWithVeto:        tallyResult.NoWithVeto.String(),
		//SystemVotingPower: tallyResult.SystemVotingPower.String(),
	}

	return
}

func GetVotes(proposalID uint64) (pVotes []document.PVote, err error) {
	cdc := types.GetCodec()

	res, err := QuerySubspace(types.KeyVotesSubspace(proposalID), "gov")
	if len(res) == 0 || err != nil {
		return pVotes, err
	}
	for i := 0; i < len(res); i++ {
		var vote types.SdkVote
		cdc.UnmarshalBinaryLengthPrefixed(res[i].Value, &vote)
		v := document.PVote{
			Voter:  vote.Voter.String(),
			Option: vote.Option.String(),
		}
		pVotes = append(pVotes, v)
	}
	return
}
