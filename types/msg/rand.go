package msg

import (
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

type DocTxMsgRequestRand struct {
	Consumer      string `bson:"consumer"`       // request address
	BlockInterval uint64 `bson:"block-interval"` // block interval after which the requested random number will be generated
}

func (doctx *DocTxMsgRequestRand) Type() string {
	return constant.TxTypeRequestRand
}

func (doctx *DocTxMsgRequestRand) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgRequestRandom)
	doctx.Consumer = msg.Consumer.String()
	doctx.BlockInterval = msg.BlockInterval
}
