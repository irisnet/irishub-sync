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

type (
	DocMsgNFTEdit struct {
		Sender string `bson:"sender"`
		ID     string `bson:"id"`
		Denom  string `bson:"denom"`
		URI    string `bson:"uri"`
		Data   string `bson:"data"`
	}
)

func (m *DocMsgNFTEdit) Type() string {
	return TxTypeNFTEdit
}

func (m *DocMsgNFTEdit) BuildMsg(v interface{}) {
	var msg types.MsgEditNFT
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.Sender = msg.Sender.String()
	m.ID = strings.ToLower(msg.ID)
	m.Denom = strings.ToLower(msg.Denom)
	m.URI = msg.URI
	m.Data = msg.Data
}

func (m *DocMsgNFTEdit) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Types = append(tx.Types, m.Type())
	tx.Addrs = append(tx.Addrs, m.Sender)
	if len(tx.Msgs) > 1 {
		return tx
	}
	tx.Type = m.Type()
	tx.From = m.Sender
	tx.To = ""
	tx.Amount = []store.Coin{}
	return tx
}
