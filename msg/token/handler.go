package token

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
	case new(types.MsgIssueToken).Type():
		var msg types.MsgIssueToken
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		docTx.From = msg.Owner.String()
		docTx.Type = constant.TxTypeAssetIssueToken
		txMsg := DocTxMsgIssueToken{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.Owner)
	case new(types.MsgEditToken).Type():
		var msg types.MsgEditToken
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		docTx.From = msg.Owner.String()
		docTx.Type = constant.TxTypeAssetEditToken
		txMsg := DocTxMsgEditToken{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.Owner)
	case new(types.MsgMintToken).Type():
		var msg types.MsgMintToken
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		docTx.From = msg.Owner.String()
		docTx.To = msg.To.String()
		docTx.Type = constant.TxTypeAssetMintToken
		txMsg := DocTxMsgMintToken{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.Owner, txMsg.To)
	case new(types.MsgTransferTokenOwner).Type():
		var msg types.MsgTransferTokenOwner
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		docTx.From = msg.SrcOwner.String()
		docTx.To = msg.DstOwner.String()
		docTx.Type = constant.TxTypeAssetTransferTokenOwner
		txMsg := DocTxMsgTransferTokenOwner{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.SrcOwner, txMsg.DstOwner)
	default:
		ok = false
	}
	return docTx, ok
}
