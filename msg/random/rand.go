package random

import (
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store"
)

type DocTxMsgRequestRand struct {
	Consumer      string      `bson:"consumer"`       // request address
	BlockInterval uint64      `bson:"block_interval"` // block interval after which the requested random number will be generated
	Oracle        bool        `bson:"oracle"`
	ServiceFeeCap store.Coins `bson:"service_fee_cap"`
}

func (doctx *DocTxMsgRequestRand) Type() string {
	return constant.TxTypeRequestRand
}

func (doctx *DocTxMsgRequestRand) BuildMsg(txMsg interface{}) {
	msg := txMsg.(types.MsgRequestRandom)
	doctx.Consumer = msg.Consumer.String()
	doctx.BlockInterval = msg.BlockInterval
	doctx.Oracle = msg.Oracle
	doctx.ServiceFeeCap = types.ParseCoins(msg.ServiceFeeCap.String())
}
