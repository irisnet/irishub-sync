package iservice

import (
	"encoding/hex"
	. "github.com/irisnet/irishub-sync/util/constant"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/store/document"
	"encoding/json"
	"github.com/irisnet/irishub-sync/types"
)

type (
	DocMsgPauseRequestContext struct {
		RequestContextID string `bson:"request_context_id" yaml:"request_context_id"`
		Consumer         string `bson:"consumer" yaml:"consumer"`
	}
)

func (m *DocMsgPauseRequestContext) Type() string {
	return MsgTypePauseRequestContext
}

func (m *DocMsgPauseRequestContext) BuildMsg(v interface{}) {
	var msg types.MsgPauseRequestContext
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.RequestContextID = hex.EncodeToString(msg.RequestContextID)
	m.Consumer = msg.Consumer.String()
}

func (m *DocMsgPauseRequestContext) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Type = m.Type()
	return tx
}
