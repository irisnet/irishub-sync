package evidence

import (
	. "github.com/irisnet/irishub-sync/msg"
	"encoding/json"
	"github.com/irisnet/irishub-sync/store/document"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/types"
	. "github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store"
)

// MsgSubmitEvidence defines an sdk.Msg type that supports submitting arbitrary
// Evidence.
type DocMsgSubmitEvidence struct {
	Submitter string `bson:"submitter"`
	Evidence  Any    `bson:"evidence"`
}

func (m *DocMsgSubmitEvidence) Type() string {
	return TxTypeSubmitEvidence
}

func (m *DocMsgSubmitEvidence) BuildMsg(v interface{}) {
	var msg types.MsgSubmitEvidence
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

}

func (m *DocMsgSubmitEvidence) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Type = m.Type()
	tx.From = m.Submitter
	tx.To = ""
	tx.Amount = []store.Coin{}
	return tx
}
