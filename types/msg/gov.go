package msg

import (
	"github.com/irisnet/irishub-sync/store"
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

type DocTxMsgSubmitProposal struct {
	Proposer       string      `bson:"proposer"`        //  Address of the proposer
	InitialDeposit store.Coins `bson:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
	Content        string      `bson:"content"`
}

type Param struct {
	Subspace string `json:"subspace" bson:"subspace"`
	Key      string `json:"key" bson:"key"`
	Value    string `json:"value" bson:"value"`
}

type Params []Param

func (doctx *DocTxMsgSubmitProposal) Type() string {
	return constant.TxTypeSubmitProposal
}

func (doctx *DocTxMsgSubmitProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgSubmitProposal)
	doctx.Content = msg.Content.String()
	doctx.Proposer = msg.Proposer.String()
	doctx.InitialDeposit = itypes.ParseCoins(msg.InitialDeposit.String())
}



// MsgVote
type DocTxMsgVote struct {
	ProposalID uint64 `bson:"proposal_id"` // ID of the proposal
	Voter      string `bson:"voter"`       //  address of the voter
	Option     string `bson:"option"`      //  option from OptionSet chosen by the voter
}

func (doctx *DocTxMsgVote) Type() string {
	return constant.TxTypeVote
}

func (doctx *DocTxMsgVote) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgVote)
	doctx.Voter = msg.Voter.String()
	doctx.Option = msg.Option.String()
	doctx.ProposalID = msg.ProposalID
}

// MsgDeposit
type DocTxMsgDeposit struct {
	ProposalID uint64      `bson:"proposal_id"` // ID of the proposal
	Depositor  string      `bson:"depositor"`   // Address of the depositor
	Amount     store.Coins `bson:"amount"`      // Coins to add to the proposal's deposit
}

func (doctx *DocTxMsgDeposit) Type() string {
	return constant.TxTypeDeposit
}

func (doctx *DocTxMsgDeposit) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgDeposit)
	doctx.Depositor = msg.Depositor.String()
	doctx.Amount = itypes.ParseCoins(msg.Amount.String())
	doctx.ProposalID = msg.ProposalID
}
