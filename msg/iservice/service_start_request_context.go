package iservice

import (
	"encoding/hex"
	. "github.com/irisnet/irishub-sync/util/constant"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/store/document"
	"encoding/json"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/store"
)

type (
	DocMsgStartRequestContext struct {
		RequestContextID string `bson:"request_context_id" yaml:"request_context_id"`
		Consumer         string `bson:"consumer" yaml:"consumer"`
	}
)

func (m *DocMsgStartRequestContext) Type() string {
	return TxTypeStartRequestContext
}

func (m *DocMsgStartRequestContext) BuildMsg(v interface{}) {
	//msg := v.(MsgStartRequestContext)
	var msg types.MsgStartRequestContext
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.RequestContextID = hex.EncodeToString(msg.RequestContextID)
	m.Consumer = msg.Consumer.String()
}

func (m *DocMsgStartRequestContext) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

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
	tx.Addrs = append(tx.Addrs, m.Consumer)
	return tx
}
