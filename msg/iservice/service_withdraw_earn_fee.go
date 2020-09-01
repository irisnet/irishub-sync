package iservice

import (
	. "github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store/document"
	"encoding/json"
	"github.com/irisnet/irishub-sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/store"
)

type (
	DocMsgWithdrawEarnedFees struct {
		Owner    string
		Provider string
	}
)

func (m *DocMsgWithdrawEarnedFees) Type() string {
	return TxTypeWithdrawEarnedFees
}

func (m *DocMsgWithdrawEarnedFees) BuildMsg(v interface{}) {
	//msg := v.(MsgWithdrawEarnedFees)
	var msg types.MsgWithdrawEarnedFees
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.Owner = msg.Owner.String()
	m.Provider = msg.Provider.String()
}

func (m *DocMsgWithdrawEarnedFees) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Addrs = append(tx.Addrs, m.Provider, m.Owner)
	tx.Types = append(tx.Types, m.Type())
	if len(tx.Msgs) > 1 {
		return tx
	}
	tx.Type = m.Type()
	if len(tx.Signers) > 0 {
		tx.From = tx.Signers[0].AddrBech32
	}
	tx.To = ""
	tx.Amount = []store.Coin{}
	return tx
}
