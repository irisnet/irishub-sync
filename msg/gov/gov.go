package gov

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

type Any struct {
	// nolint
	TypeUrl string `bson:"type_url"`
	// Must be a valid serialized protocol buffer of the above specified type.
	Value string `bson:"value"`

	//cachedValue interface{} `bson:"cached_value"`

	//compat anyCompat
}

//type anyCompat struct {
//	aminoBz []byte
//	jsonBz  []byte
//	err     error
//}

type DocTxMsgSubmitProposal struct {
	Proposer       string      `bson:"proposer"`        //  Address of the proposer
	InitialDeposit store.Coins `bson:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
	Content        Any         `bson:"content"`
}

func (doctx *DocTxMsgSubmitProposal) Type() string {
	return constant.TxTypeSubmitProposal
}

func (doctx *DocTxMsgSubmitProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(types.MsgSubmitProposal)
	doctx.Content = Any{
		TypeUrl: msg.Content.GetTypeUrl(),
		//cachedValue:msg.Content.GetCachedValue(),
		Value: string(msg.Content.GetValue()),
	}
	doctx.Proposer = msg.Proposer.String()
	doctx.InitialDeposit = types.ParseCoins(msg.InitialDeposit.String())
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
	msg := txMsg.(types.MsgVote)
	doctx.Voter = msg.Voter.String()
	doctx.Option = msg.Option.String()
	doctx.ProposalID = msg.ProposalId
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
	msg := txMsg.(types.MsgDeposit)
	doctx.Depositor = msg.Depositor.String()
	doctx.Amount = types.ParseCoins(msg.Amount.String())
	doctx.ProposalID = msg.ProposalId
}
