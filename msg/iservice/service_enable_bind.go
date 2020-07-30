package iservice

import (
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	. "github.com/irisnet/irishub-sync/util/constant"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
	"github.com/irisnet/irishub-sync/store"
)

type (
	DocMsgEnableServiceBinding struct {
		ServiceName string      `bson:"service_name" yaml:"service_name"`
		Provider    string      `bson:"provider" yaml:"provider"`
		Deposit     store.Coins `bson:"deposit" yaml:"deposit"`
		Owner       string      `bson:"owner" yaml:"owner"`
	}
)

func (m *DocMsgEnableServiceBinding) Type() string {
	return TxTypeEnableServiceBinding
}

func (m *DocMsgEnableServiceBinding) BuildMsg(v interface{}) {
	var msg types.MsgEnableServiceBinding
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	var coins store.Coins
	for _, one := range msg.Deposit {
		coins = append(coins, types.ParseCoin(one.Amount.String()))
	}

	m.ServiceName = msg.ServiceName
	m.Provider = msg.Provider.String()
	m.Deposit = coins
	m.Owner = msg.Owner.String()
}

func (m *DocMsgEnableServiceBinding) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

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
	tx.Amount = m.Deposit
	tx.Addrs = append(tx.Addrs, m.Provider, m.Owner)
	return tx
}
