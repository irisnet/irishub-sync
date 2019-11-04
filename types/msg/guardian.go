package msg

import (
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

type DocTxMsgAddProfiler struct {
	AddGuardian
}

func (doctx *DocTxMsgAddProfiler) Type() string {
	return constant.TxTypeAddProfiler
}

func (doctx *DocTxMsgAddProfiler) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgAddProfiler)
	doctx.Address = msg.Address.String()
	doctx.AddedBy = msg.AddedBy.String()
	doctx.Description = msg.Description
}

type DocTxMsgAddTrustee struct {
	AddGuardian
}

func (doctx *DocTxMsgAddTrustee) Type() string {
	return constant.TxTypeAddTrustee
}

func (doctx *DocTxMsgAddTrustee) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgAddTrustee)
	doctx.Address = msg.Address.String()
	doctx.AddedBy = msg.AddedBy.String()
	doctx.Description = msg.Description
}

type AddGuardian struct {
	Description string `bson:"description"`
	Address     string `bson:"address"`  // address added
	AddedBy     string `bson:"added_by"` // address that initiated the tx
}

type DocTxMsgDeleteProfiler struct {
	DeleteGuardian
}

func (doctx *DocTxMsgDeleteProfiler) Type() string {
	return constant.TxTypeDeleteProfiler
}

func (doctx *DocTxMsgDeleteProfiler) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgDeleteProfiler)
	doctx.Address = msg.Address.String()
	doctx.DeletedBy = msg.DeletedBy.String()
}

type DocTxMsgDeleteTrustee struct {
	DeleteGuardian
}

func (doctx *DocTxMsgDeleteTrustee) Type() string {
	return constant.TxTypeDeleteTrustee
}

func (doctx *DocTxMsgDeleteTrustee) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgDeleteTrustee)
	doctx.Address = msg.Address.String()
	doctx.DeletedBy = msg.DeletedBy.String()
}

type DeleteGuardian struct {
	Address   string `bson:"address"`    // address deleted
	DeletedBy string `bson:"deleted_by"` // address that initiated the tx
}
