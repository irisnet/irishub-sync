package types

import (
	"github.com/irisnet/irishub-sync/util/constant"
)

type DocTxMsgSubmitTokenAdditionProposal struct {
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
	msg := txMsg.(MsgSubmitTokenAdditionProposal)
	doctx.Symbol = msg.Symbol
	doctx.MinUnitAlias = msg.MinUnitAlias
	doctx.CanonicalSymbol = msg.CanonicalSymbol
	doctx.Name = msg.Name
	doctx.Decimal = msg.Decimal
	doctx.InitialSupply = msg.InitialSupply
}

type DocTxMsgSetMemoRegexp struct {
	Owner      string `bson:"owner"`
	MemoRegexp string `bson:"memo_regexp"`
}

func (doctx *DocTxMsgSetMemoRegexp) Type() string {
	return constant.TxTypeSetMemoRegexp
}

func (doctx *DocTxMsgSetMemoRegexp) BuildMsg(txMsg interface{}) {
	msg := txMsg.(MsgSetMemoRegexp)
	doctx.MemoRegexp = msg.MemoRegexp
	doctx.Owner = msg.Owner.String()
}

type DocTxMsgRequestRand struct {
	Consumer      string `bson:"consumer"`       // request address
	BlockInterval uint64 `bson:"block-interval"` // block interval after which the requested random number will be generated
}

func (doctx *DocTxMsgRequestRand) Type() string {
	return constant.TxTypeRequestRand
}

func (doctx *DocTxMsgRequestRand) BuildMsg(txMsg interface{}) {
	msg := txMsg.(MsgRequestRand)
	doctx.Consumer = msg.Consumer.String()
	doctx.BlockInterval = msg.BlockInterval
}
