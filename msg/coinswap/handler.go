package coinswap

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
	case new(types.MsgAddLiquidity).Type():
		var msg types.MsgAddLiquidity
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgAddLiquidity{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.Sender)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.Sender.String()
		docTx.To = ""
		docTx.Amount = types.ParseCoins(msg.MaxToken.String())
		docTx.Type = constant.TxTypeAddLiquidity
	case new(types.MsgRemoveLiquidity).Type():
		var msg types.MsgRemoveLiquidity
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgRemoveLiquidity{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.Sender)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.Sender.String()
		docTx.To = ""
		docTx.Amount = types.ParseCoins(msg.WithdrawLiquidity.String())
		docTx.Type = constant.TxTypeRemoveLiquidity
	case new(types.MsgSwapOrder).Type():
		var msg types.MsgSwapOrder
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgSwapOrder{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.Input.Address, txMsg.Output.Address)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.Input.Address.String()
		docTx.To = msg.Output.Address.String()
		docTx.Amount = types.ParseCoins(msg.Input.Coin.String())
		docTx.Type = constant.TxTypeSwapOrder
	default:
		ok = false
	}
	return docTx, ok
}
