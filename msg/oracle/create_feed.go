package oracle

import (
	"github.com/irisnet/irishub-sync/store"
	"encoding/json"
	"github.com/irisnet/irishub-sync/store/document"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/types"
	. "github.com/irisnet/irishub-sync/util/constant"
)

type DocMsgCreateFeed struct {
	FeedName          string      `bson:"feed_name" yaml:"feed_name"`
	LatestHistory     uint64      `bson:"latest_history" yaml:"latest_history"`
	Description       string      `bson:"description"`
	Creator           string      `bson:"creator"`
	ServiceName       string      `bson:"service_name" yaml:"service_name"`
	Providers         []string    `bson:"providers"`
	Input             string      `bson:"input"`
	Timeout           int64       `bson:"timeout"`
	ServiceFeeCap     store.Coins `bson:"service_fee_cap" yaml:"service_fee_cap"`
	RepeatedFrequency uint64      `bson:"repeated_frequency" yaml:"repeated_frequency"`
	AggregateFunc     string      `bson:"aggregate_func" yaml:"aggregate_func"`
	ValueJsonPath     string      `bson:"value_json_path" yaml:"value_json_path"`
	ResponseThreshold uint32      `bson:"response_threshold" yaml:"response_threshold"`
}

func (m *DocMsgCreateFeed) Type() string {
	return TxTypeCreateFeed
}

func (m *DocMsgCreateFeed) BuildMsg(v interface{}) {
	var msg types.MsgCreateFeed
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

}

func (m *DocMsgCreateFeed) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Type = m.Type()
	return tx
}
