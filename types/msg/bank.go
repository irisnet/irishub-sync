package msg

import (
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
