package oracle

import (
	"github.com/irisnet/irishub-sync/store"
	"encoding/json"
	"github.com/irisnet/irishub-sync/store/document"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/types"
	. "github.com/irisnet/irishub-sync/util/constant"
)

type DocMsgEditFeed struct {
	FeedName          string      `bson:"feed_name" yaml:"feed_name"`
	LatestHistory     uint64      `bson:"latest_history" yaml:"latest_history"`
	Description       string      `bson:"description"`
	Creator           string      `bson:"creator"`
	Providers         []string    `bson:"providers"`
	Timeout           int64       `bson:"timeout"`
	ServiceFeeCap     store.Coins `bson:"service_fee_cap" yaml:"service_fee_cap"`
	RepeatedFrequency uint64      `bson:"repeated_frequency" yaml:"repeated_frequency"`
	ResponseThreshold uint32      `bson:"response_threshold" yaml:"response_threshold"`
}

func (m *DocMsgEditFeed) Type() string {
	return TxTypeEditFeed
}

func (m *DocMsgEditFeed) BuildMsg(v interface{}) {
	var msg types.MsgCreateFeed
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.FeedName = msg.FeedName
	m.LatestHistory = msg.LatestHistory
	m.Description = msg.Description
	m.Creator = msg.Creator.String()
	for _, val := range msg.GetProviders() {
		m.Providers = append(m.Providers, val.String())
	}
	m.Timeout = msg.Timeout
	m.ServiceFeeCap = types.ParseCoins(msg.ServiceFeeCap.String())
	m.RepeatedFrequency = msg.RepeatedFrequency
	m.ResponseThreshold = msg.ResponseThreshold
}

func (m *DocMsgEditFeed) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Addrs = append(tx.Addrs, m.Providers...)
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
