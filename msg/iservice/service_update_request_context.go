package iservice

import (
	"encoding/hex"
	. "github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store/document"
	"encoding/json"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	DocMsgUpdateRequestContext struct {
		RequestContextID  string      `bson:"request_context_id" yaml:"request_context_id"`
		Providers         []string    `bson:"providers" yaml:"providers"`
		Consumer          string      `bson:"consumer" yaml:"consumer"`
		ServiceFeeCap     store.Coins `bson:"service_fee_cap" yaml:"service_fee_cap"`
		Timeout           int64       `bson:"timeout" yaml:"timeout"`
		RepeatedFrequency uint64      `bson:"repeated_frequency" yaml:"repeated_frequency"`
		RepeatedTotal     int64       `bson:"repeated_total" yaml:"repeated_total"`
	}
)

func (m *DocMsgUpdateRequestContext) Type() string {
	return MsgTypeUpdateRequestContext
}

func (m *DocMsgUpdateRequestContext) BuildMsg(v interface{}) {
	//msg := v.(MsgUpdateRequestContext)
	var msg types.MsgUpdateRequestContext
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	loadProviders := func() (ret []string) {
		for _, one := range msg.Providers {
			ret = append(ret, one.String())
		}
		return
	}

	var coins store.Coins
	for _, one := range msg.ServiceFeeCap {
		coins = append(coins, types.ParseCoin(one.String()))
	}

	m.RequestContextID = hex.EncodeToString(msg.RequestContextID)
	m.Providers = loadProviders()
	m.Consumer = msg.Consumer.String()
	m.ServiceFeeCap = coins
	m.Timeout = msg.Timeout
	m.RepeatedFrequency = msg.RepeatedFrequency
	m.RepeatedTotal = msg.RepeatedTotal
}

func (m *DocMsgUpdateRequestContext) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

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
