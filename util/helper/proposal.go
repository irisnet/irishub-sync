package helper

import (
	"errors"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

var ProposalStatusName = map[int32]string{
	0: constant.StatusUnspecified,
	1: constant.StatusDepositPeriod,
	2: constant.StatusVotingPeriod,
	3: constant.StatusPassed,
	4: constant.StatusRejected,
	5: constant.StatusFailed,
}

func GetProposal(proposalID uint64) (proposal document.Proposal, err error) {

	res, err := Query(types.KeyProposal(proposalID), "gov", constant.StoreDefaultEndPath)
	if len(res) == 0 || err != nil {
		return proposal, errors.New("no data")
	}

	var propo types.Proposal
	var txtpropo types.TextProposal
	propo.Unmarshal(res) //TODO
	txtpropo.Unmarshal(propo.Content.Value)
	proposal.ProposalId = proposalID
	proposal.Title = txtpropo.Title
	proposal.Description = txtpropo.Description
	if stat, ok := ProposalStatusName[int32(propo.Status)]; ok {
		proposal.Status = stat
	} else {
		proposal.Status = propo.Status.String()
	}

	proposal.SubmitTime = propo.SubmitTime
	proposal.VotingStartTime = propo.VotingStartTime
	proposal.VotingEndTime = propo.VotingEndTime
	proposal.DepositEndTime = propo.DepositEndTime
	proposal.TotalDeposit = types.ParseCoins(propo.TotalDeposit.String())
	proposal.Votes = []document.PVote{}

	tallyResult := propo.FinalTallyResult
	proposal.TallyResult = document.PTallyResult{
		Yes:        tallyResult.Yes.String(),
		Abstain:    tallyResult.Abstain.String(),
		No:         tallyResult.No.String(),
		NoWithVeto: tallyResult.NoWithVeto.String(),
	}

	return
}
