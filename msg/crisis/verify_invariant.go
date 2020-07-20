package crisis

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/store/document"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/types"
	. "github.com/irisnet/irishub-sync/util/constant"
)

type DocMsgVerifyInvariant struct {
	Sender              string `bson:"sender"`
	InvariantModuleName string `bson:"invariant_module_name" yaml:"invariant_module_name"`
	InvariantRoute      string `bson:"invariant_route" yaml:"invariant_route"`
}

func (m *DocMsgVerifyInvariant) Type() string {
	return TxTypeVerifyInvariant
}

func (m *DocMsgVerifyInvariant) BuildMsg(v interface{}) {
	var msg types.MsgVerifyInvariant
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

}

func (m *DocMsgVerifyInvariant) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Type = m.Type()
	return tx
}
