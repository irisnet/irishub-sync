package nft

import (
	. "github.com/irisnet/irishub-sync/util/constant"
	"strings"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
	"github.com/irisnet/irishub-sync/store"
)

type DocMsgIssueDenom struct {
	Sender string `bson:"sender"`
	ID     string `bson:"id"`
	Name   string `bson:"name"`
	Schema string `bson:"schema"`
}

func (m *DocMsgIssueDenom) Type() string {
	return TxTypeIssueDenom
}

func (m *DocMsgIssueDenom) BuildMsg(v interface{}) {
	var msg types.MsgIssueDenom
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.Sender = msg.Sender.String()
	m.Schema = msg.Schema
	m.Name = msg.Name
	m.ID = strings.ToLower(msg.ID)
}

func (m *DocMsgIssueDenom) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Addrs = append(tx.Addrs, m.Sender)
	tx.Types = append(tx.Types, m.Type())
	if len(tx.Msgs) > 1 {
		return tx
	}
	tx.From = m.Sender
	tx.To = ""
	tx.Amount = []store.Coin{}
	tx.Type = m.Type()
	return tx
}
