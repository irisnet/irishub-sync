package iservice

import (
	"encoding/hex"
	. "github.com/irisnet/irishub-sync/util/constant"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/store/document"
	"encoding/json"
	"github.com/irismod/service/types"
)

type (
	DocMsgServiceResponse struct {
		RequestID string `bson:"request_id" yaml:"request_id"`
		Provider  string `bson:"provider" yaml:"provider"`
		Output    string `bson:"output" yaml:"output"`
		Result    string `bson:"result"`
	}
)

func (m *DocMsgServiceResponse) Type() string {
	return MsgTypeRespondService
}

func (m *DocMsgServiceResponse) BuildMsg(msg interface{}) {
	//v := msg.(MsgRespondService)
	var v types.MsgRespondService
	data, _ := json.Marshal(msg)
	json.Unmarshal(data, &v)

	m.RequestID = hex.EncodeToString(v.RequestID)
	m.Provider = v.Provider.String()
	//m.Output = hex.EncodeToString(v.Output)
	m.Output = v.Output
	m.Result = v.Result
}

func (m *DocMsgServiceResponse) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Type = m.Type()
	return tx
}
