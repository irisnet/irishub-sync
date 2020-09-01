package bank

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
	case new(types.MsgTransfer).Type():
		var msg types.MsgTransfer
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgSend{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.ToAddress, txMsg.FromAddress)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.FromAddress.String()
		docTx.To = msg.ToAddress.String()
		docTx.Amount = types.ParseCoins(msg.Amount.String())
		docTx.Type = constant.TxTypeTransfer
	case new(types.MsgMultiSend).Type():
		var msg types.MsgMultiSend
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)
		txMsg := DocMsgMultiSend{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.Temp...)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.Type = constant.TxTypeMultiSend
	default:
		ok = false
	}
	return docTx, ok
}
