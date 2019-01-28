package types

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/constant"
)

type SubmitProposal struct {
	Title          string      `json:"title"`          //  Title of the proposal
	Description    string      `json:"description"`    //  Description of the proposal
	Proposer       string      `json:"proposer"`       //  Address of the proposer
	InitialDeposit store.Coins `json:"initialDeposit"` //  Initial deposit paid by sender. Must be strictly positive.
	ProposalType   string      `json:"proposalType"`   //  Initial deposit paid by sender. Must be strictly positive.
	Params         Params      `json:"params"`
}

type SubmitSoftwareUpgradeProposal struct {
	SubmitProposal
	Version      uint64 `json:"version"`
	Software     string `json:"software"`
	SwitchHeight uint64 `json:"switch_height"`
	Threshold    string `json:"threshold"`
}

type Param struct {
	Subspace string `json:"subspace"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type Params []Param

func NewSubmitProposal(msg MsgSubmitProposal) SubmitProposal {
	var params Params
	for _, p := range msg.Params {
		params = append(params, Param{
			Subspace: p.Subspace,
			Key:      p.Key,
			Value:    p.Value,
		})
	}
	return SubmitProposal{
		Title:          msg.Title,
		Description:    msg.Description,
		ProposalType:   msg.ProposalType.String(),
		Proposer:       msg.Proposer.String(),
		InitialDeposit: ParseCoins(msg.InitialDeposit.String()),
		Params:         params,
	}
}
func NewSubmitSoftwareUpgradeProposal(msg MsgSubmitSoftwareUpgradeProposal) SubmitSoftwareUpgradeProposal {
	submitProposal := NewSubmitProposal(msg.MsgSubmitProposal)
	return SubmitSoftwareUpgradeProposal{
		SubmitProposal: submitProposal,
		Version:        msg.Version,
		Software:       msg.Software,
		SwitchHeight:   msg.SwitchHeight,
		Threshold:      msg.Threshold.String(),
	}
}

func (s SubmitProposal) Type() string {
	return constant.TxTypeSubmitProposal
}

func (s SubmitProposal) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func (s SubmitSoftwareUpgradeProposal) Type() string {
	return constant.TxTypeSubmitProposal
}

func (s SubmitSoftwareUpgradeProposal) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func UnmarshalSubmitProposal(str string) (submitProposal SubmitProposal) {
	json.Unmarshal([]byte(str), &submitProposal)
	return
}

type Vote struct {
	ProposalID uint64 `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     string `json:"option"`
}

func NewVote(v MsgVote) Vote {
	return Vote{
		ProposalID: v.ProposalID,
		Voter:      v.Voter.String(),
		Option:     v.Option.String(),
	}
}

func (s Vote) Type() string {
	return constant.TxTypeVote
}

func (s Vote) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func UnmarshalVote(str string) (vote Vote) {
	json.Unmarshal([]byte(str), &vote)
	return
}

type Deposit struct {
	ProposalID uint64      `json:"proposal_id"` // ID of the proposal
	Depositer  string      `json:"depositer"`   // Address of the depositer
	Amount     store.Coins `json:"amount"`      // Coins to add to the proposal's deposit
}

func NewDeposit(deposit MsgDeposit) Deposit {
	return Deposit{
		ProposalID: deposit.ProposalID,
		Depositer:  deposit.Depositor.String(),
		Amount:     ParseCoins(deposit.Amount.String()),
	}
}

func (s Deposit) Type() string {
	return constant.TxTypeDeposit
}

func (s Deposit) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func UnmarshalDeposit(str string) (deposit Deposit) {
	json.Unmarshal([]byte(str), &deposit)
	return
}
