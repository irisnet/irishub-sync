package types

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/constant"
)

type SubmitProposal struct {
	Title          string      `json:"title"`          //  Title of the proposal
	Description    string      `json:"description"`    //  Description of the proposal
	ProposalType   string      `json:"proposalType"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       string      `json:"proposer"`       //  Address of the proposer
	InitialDeposit store.Coins `json:"initialDeposit"` //  Initial deposit paid by sender. Must be strictly positive.
	Params         Param       `json:"params"`
}

type Param struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
	Op    string `bson:"op"`
}

func NewSubmitProposal(msg MsgSubmitProposal) SubmitProposal {
	p := Param{
		Key:   msg.Param.Key,
		Value: msg.Param.Value,
		Op:    msg.Param.Op,
	}
	return SubmitProposal{
		Title:          msg.Title,
		Description:    msg.Description,
		ProposalType:   msg.ProposalType.String(),
		Proposer:       msg.Proposer.String(),
		InitialDeposit: BuildCoins(msg.InitialDeposit),
		Params:         p,
	}
}

func (s SubmitProposal) Type() string {
	return constant.TxTypeSubmitProposal
}

func (s SubmitProposal) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func UnmarshalSubmitProposal(str string) (submitProposal SubmitProposal) {
	json.Unmarshal([]byte(str), &submitProposal)
	return
}
