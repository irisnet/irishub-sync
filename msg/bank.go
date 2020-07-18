package msg

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
