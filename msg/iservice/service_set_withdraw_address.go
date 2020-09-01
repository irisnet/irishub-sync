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
	DocMsgServiceSetWithdrawAddress struct {
		Owner           string `bson:"owner" yaml:"owner"`
		WithdrawAddress string `bson:"withdraw_address" yaml:"withdraw_address"`
	}
)

func (m *DocMsgServiceSetWithdrawAddress) Type() string {
	return TxTypeServiceSetWithdrawAddress
}

func (m *DocMsgServiceSetWithdrawAddress) BuildMsg(v interface{}) {
	var msg types.MsgSetWithdrawFeesAddress
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.Owner = msg.Owner.String()
	m.WithdrawAddress = msg.WithdrawAddress.String()
}

func (m *DocMsgServiceSetWithdrawAddress) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) (*document.CommonTx, bool) {

	m.BuildMsg(msgData)
	if m.Owner == "" {
		return tx, false
	}
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Addrs = append(tx.Addrs, m.Owner)
	tx.Types = append(tx.Types, m.Type())
	if len(tx.Msgs) > 1 {
		return tx, true
	}
	tx.Type = m.Type()
	if len(tx.Signers) > 0 {
		tx.From = tx.Signers[0].AddrBech32
	}
	tx.To = ""
	tx.Amount = []store.Coin{}
	tx.Addrs = append(tx.Addrs, m.Owner)
	return tx, true
}
