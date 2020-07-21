package iservice

import (
	. "github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
	"github.com/irisnet/irishub-sync/store"
)

type (
	DocMsgDisableServiceBinding struct {
		ServiceName string `bson:"service_name" yaml:"service_name"`
		Provider    string `bson:"provider" yaml:"provider"`
		Owner       string `bson:"owner" yaml:"owner"`
	}
)

func (m *DocMsgDisableServiceBinding) Type() string {
	return MsgTypeDisableServiceBinding
}

func (m *DocMsgDisableServiceBinding) BuildMsg(v interface{}) {
	var msg types.MsgDisableServiceBinding
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)
	m.ServiceName = msg.ServiceName
	m.Provider = msg.Provider.String()
	m.Owner = msg.Owner.String()
}

func (m *DocMsgDisableServiceBinding) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

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
