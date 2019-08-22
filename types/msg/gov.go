package msg

import (
	"github.com/irisnet/irishub-sync/store"
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

type DocTxMsgSubmitProposal struct {
	Title          string        `bson:"title"`          //  Title of the proposal
	Description    string        `bson:"description"`    //  Description of the proposal
	Proposer       string        `bson:"proposer"`       //  Address of the proposer
	InitialDeposit store.Coins   `bson:"initialDeposit"` //  Initial deposit paid by sender. Must be strictly positive.
	ProposalType   string        `bson:"proposalType"`   //  Initial deposit paid by sender. Must be strictly positive.
	Params         itypes.Params `bson:"params"`
}

type DocTxMsgSubmitTokenAdditionProposal struct {
	DocTxMsgSubmitProposal
	Symbol          string `bson:"symbol"`
	CanonicalSymbol string `bson:"canonical_symbol"`
	Name            string `bson:"name"`
	Decimal         uint8  `bson:"decimal"`
	MinUnitAlias    string `bson:"min_unit_alias"`
	InitialSupply   uint64 `bson:"initial_supply"`
}

func (doctx *DocTxMsgSubmitTokenAdditionProposal) Type() string {
	return constant.TxMsgTypeSubmitTokenAdditionProposal
}

func (doctx *DocTxMsgSubmitTokenAdditionProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgSubmitTokenAdditionProposal)
	var params itypes.Params
	for _, p := range msg.Params {
		params = append(params, itypes.Param{
			Subspace: p.Subspace,
			Key:      p.Key,
			Value:    p.Value,
		})
	}
	doctx.Title = msg.Title
	doctx.Description = msg.Description
	doctx.ProposalType = msg.ProposalType.String()
	doctx.Proposer = msg.Proposer.String()
	doctx.Params = params
	doctx.InitialDeposit = itypes.ParseCoins(msg.InitialDeposit.String())
	doctx.Symbol = msg.Symbol
	doctx.MinUnitAlias = msg.MinUnitAlias
	doctx.CanonicalSymbol = msg.CanonicalSymbol
	doctx.Name = msg.Name
	doctx.Decimal = msg.Decimal
	doctx.InitialSupply = msg.InitialSupply
}
