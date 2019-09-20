package msg

import (
	"github.com/irisnet/irishub-sync/util/constant"
	itypes "github.com/irisnet/irishub-sync/types"
)

// MsgDelegate - struct for bonding transactions
type DocTxMsgBeginRedelegate struct {
	DelegatorAddr    string `bson:"delegator_addr"`
	ValidatorSrcAddr string `bson:"validator_src_addr"`
	ValidatorDstAddr string `bson:"validator_dst_addr"`
	SharesAmount     string `bson:"shares_amount"`
}

func (doctx *DocTxMsgBeginRedelegate) Type() string {
	return constant.TxTypeBeginRedelegate
}

func (doctx *DocTxMsgBeginRedelegate) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgBeginRedelegate)
	doctx.DelegatorAddr = msg.DelegatorAddr.String()
	doctx.ValidatorSrcAddr = msg.ValidatorSrcAddr.String()
	doctx.ValidatorDstAddr = msg.ValidatorDstAddr.String()
	doctx.SharesAmount = msg.SharesAmount.String()
}
