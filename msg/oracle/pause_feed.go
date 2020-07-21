package oracle

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store"
)

type DocMsgPauseFeed struct {
	FeedName string `bson:"feed_name" yaml:"feed_name"`
	Creator  string `bson:"creator"`
}

func (m *DocMsgPauseFeed) Type() string {
	return TxTypePauseFeed
}

func (m *DocMsgPauseFeed) BuildMsg(v interface{}) {
	var msg types.MsgCreateFeed
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

}

func (m *DocMsgPauseFeed) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Type = m.Type()
	tx.From = m.Creator
	tx.To = ""
	tx.Amount = []store.Coin{}
	return tx
}
