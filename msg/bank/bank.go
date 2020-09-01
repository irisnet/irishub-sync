package bank

import (
	"github.com/irisnet/irishub-sync/store"
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

// MsgSend - high level transaction of the coin module
type DocTxMsgSend struct {
	FromAddress string      `bson:"from_address"`
	ToAddress   string      `bson:"to_address"`
	Amount      store.Coins `bson:"amount"`
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
	doctx.FromAddress = msg.FromAddress.String()
	doctx.ToAddress = msg.ToAddress.String()
	doctx.Amount = itypes.ParseCoins(msg.Amount.String())
}

type (
	DocMsgMultiSend struct {
		Inputs  []Item   `bson:"inputs"`
		Outputs []Item   `bson:"outputs"`
		Temp    []string `bson:"-"`
	}
	Item struct {
		Address string      `bson:"address"`
		Coins   store.Coins `bson:"coins"`
	}
)

func (m *DocMsgMultiSend) Type() string {
	return constant.TxTypeMultiSend
}

func (m *DocMsgMultiSend) BuildMsg(v interface{}) {
	msg := v.(itypes.MsgMultiSend)
	for _, one := range msg.Inputs {
		m.Inputs = append(m.Inputs, Item{Address: one.Address.String(), Coins: itypes.ParseCoins(one.Coins.String())})
		m.Temp = append(m.Temp, one.Address.String())
	}
	for _, one := range msg.Outputs {
		m.Outputs = append(m.Outputs, Item{Address: one.Address.String(), Coins: itypes.ParseCoins(one.Coins.String())})
		m.Temp = append(m.Temp, one.Address.String())
	}

}
