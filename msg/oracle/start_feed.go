package oracle

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store"
)

type DocMsgStartFeed struct {
	FeedName string `bson:"feed_name" yaml:"feed_name"`
	Creator  string `bson:"creator"`
}

func (m *DocMsgStartFeed) Type() string {
	return TxTypeStartFeed
}

func (m *DocMsgStartFeed) BuildMsg(v interface{}) {
	var msg types.MsgStartFeed
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.FeedName = msg.FeedName
	m.Creator = msg.Creator.String()
}

func (m *DocMsgStartFeed) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Addrs = append(tx.Addrs, m.Creator)
	tx.Types = append(tx.Types, m.Type())
	if len(tx.Msgs) > 1 {
		return tx
	}
	tx.Type = m.Type()
	tx.From = m.Creator
	tx.To = ""
	tx.Amount = []store.Coin{}
	return tx
}
