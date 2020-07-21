package iservice

import (
	. "github.com/irisnet/irishub-sync/util/constant"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/store/document"
	"encoding/json"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/store"
)

type (
	DocMsgSetWithdrawAddress struct {
		DelegatorAddress string `bson:"delegator_address" yaml:"delegator_address"`
		WithdrawAddress  string `bson:"withdraw_address" yaml:"withdraw_address"`
	}
)

func (m *DocMsgSetWithdrawAddress) Type() string {
	return MsgTypeSetWithdrawAddress
}

func (m *DocMsgSetWithdrawAddress) BuildMsg(v interface{}) {
	//msg := v.(MsgSetWithdrawAddress)
	var msg types.MsgSetWithdrawAddress
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.DelegatorAddress = msg.DelegatorAddress.String()
	m.WithdrawAddress = msg.WithdrawAddress.String()
}

func (m *DocMsgSetWithdrawAddress) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Type = m.Type()
	if len(tx.Signers) > 0 {
		tx.From = tx.Signers[0].AddrBech32
	}
	tx.To = ""
	tx.Amount = []store.Coin{}
	return tx
}
