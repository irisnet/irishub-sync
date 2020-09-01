package guardian

import (
	"github.com/irisnet/irishub-sync/store/document"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/types"
	"encoding/json"
	"github.com/irisnet/irishub-sync/util/constant"
)

func HandleTxMsg(msgData sdk.Msg, docTx *document.CommonTx) (*document.CommonTx, bool) {
	ok := true
	switch msgData.Type() {
	case new(types.MsgAddProfiler).Type():
		var msg types.MsgAddProfiler
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)
		txMsg := DocTxMsgAddProfiler{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.Address, txMsg.AddedBy)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.AddGuardian.AddedBy.String()
		docTx.To = msg.AddGuardian.Address.String()
		docTx.Type = constant.TxTypeAddProfiler

	case new(types.MsgAddTrustee).Type():
		var msg types.MsgAddTrustee
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)


		txMsg := DocTxMsgAddTrustee{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.Address, txMsg.AddedBy)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.AddGuardian.AddedBy.String()
		docTx.To = msg.AddGuardian.Address.String()
		docTx.Type = constant.TxTypeAddTrustee

	case new(types.MsgDeleteTrustee).Type():
		var msg types.MsgDeleteTrustee
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)


		txMsg := DocTxMsgDeleteTrustee{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.DeletedBy, txMsg.Address)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.DeleteGuardian.DeletedBy.String()
		docTx.To = msg.DeleteGuardian.Address.String()
		docTx.Type = constant.TxTypeDeleteTrustee

	case new(types.MsgDeleteProfiler).Type():
		var msg types.MsgDeleteProfiler
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgDeleteProfiler{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.DeletedBy, txMsg.Address)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.DeleteGuardian.DeletedBy.String()
		docTx.To = msg.DeleteGuardian.Address.String()
		docTx.Type = constant.TxTypeDeleteProfiler
	default:
		ok = false
	}
	return docTx, ok
}

