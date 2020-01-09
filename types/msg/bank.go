package msg

import (
	"github.com/irisnet/irishub-sync/store"
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

type DocTxMsgSetMemoRegexp struct {
	Owner      string `bson:"owner"`
	MemoRegexp string `bson:"memo_regexp"`
}

func (doctx *DocTxMsgSetMemoRegexp) Type() string {
	return constant.TxTypeSetMemoRegexp
}

func (doctx *DocTxMsgSetMemoRegexp) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgSetMemoRegexp)
	doctx.MemoRegexp = msg.MemoRegexp
	doctx.Owner = msg.Owner.String()
}

// MsgBurn - high level transaction of the coin module
type DocTxMsgBurn struct {
	Owner string      `bson:"owner"`
	Coins store.Coins `bson:"coins"`
}

func (doctx *DocTxMsgBurn) Type() string {
	return constant.TxTypeBurn
}

func (doctx *DocTxMsgBurn) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgBurn)
	doctx.Owner = msg.Owner.String()
	doctx.Coins = itypes.ParseCoins(msg.Coins.String())
}

// MsgSend - high level transaction of the coin module
type DocTxMsgSend struct {
	Inputs  []Data `bson:"inputs"`
	Outputs []Data `bson:"outputs"`
}

// Transaction
type Data struct {
	Address string      `bson:"address"`
	Coins   store.Coins `bson:"coins"`
}

func (doctx *DocTxMsgSend) Type() string {
	return constant.TxTypeTransfer
}

func (doctx *DocTxMsgSend) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgTransfer)
	doctx.Inputs = append(doctx.Inputs, Data{
		Address: msg.Inputs[0].Address.String(),
		Coins:   itypes.ParseCoins(msg.Inputs[0].Coins.String()),
	})
	doctx.Outputs = append(doctx.Outputs, Data{
		Address: msg.Outputs[0].Address.String(),
		Coins:   itypes.ParseCoins(msg.Outputs[0].Coins.String()),
	})
}
