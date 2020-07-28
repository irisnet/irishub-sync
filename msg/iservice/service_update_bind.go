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
	DocMsgUpdateServiceBinding struct {
		ServiceName string      `bson:"service_name" yaml:"service_name"`
		Provider    string      `bson:"provider" yaml:"provider"`
		Deposit     store.Coins `bson:"deposit" yaml:"deposit"`
		Pricing     string      `bson:"pricing" yaml:"pricing"`
		QoS         uint64      `bson:"qos" yaml:"qos"`
		Owner       string      `bson:"owner" yaml:"owner"`
	}
)

func (m *DocMsgUpdateServiceBinding) Type() string {
	return TxTypeUpdateServiceBinding
}

func (m *DocMsgUpdateServiceBinding) BuildMsg(v interface{}) {
	//msg := v.(MsgUpdateServiceBinding)
	var msg types.MsgUpdateServiceBinding
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	var coins store.Coins
	for _, one := range msg.Deposit {
		coins = append(coins, types.ParseCoin(one.String()))
	}

	m.ServiceName = msg.ServiceName
	m.Provider = msg.Provider.String()
	m.Deposit = coins
	m.Pricing = msg.Pricing
	m.QoS = msg.QoS
	m.Owner = msg.Owner.String()
}

func (m *DocMsgUpdateServiceBinding) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

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
	tx.Addrs = append(tx.Addrs, m.Provider, m.Owner)
	return tx
}
