package staking

import (
	"github.com/irisnet/irishub-sync/store/document"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/types"
	"encoding/json"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store"
)

func HandleTxMsg(msgData sdk.Msg, docTx *document.CommonTx) (*document.CommonTx, bool) {
	ok := true
	switch msgData.Type() {
	case new(types.MsgStakeCreate).Type():
		var msg types.MsgStakeCreate
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgStakeCreate{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.DelegatorAddr, txMsg.ValidatorAddr)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.DelegatorAddress.String()
		docTx.To = msg.ValidatorAddress.String()
		docTx.Amount = []store.Coin{types.ParseCoin(msg.Value.String())}
		docTx.Type = constant.TxTypeStakeCreateValidator
	case new(types.MsgStakeEdit).Type():
		var msg types.MsgStakeEdit
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgStakeEdit{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.ValidatorAddr)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.ValidatorAddress.String()
		docTx.To = ""
		docTx.Amount = []store.Coin{}
		docTx.Type = constant.TxTypeStakeEditValidator

	case new(types.MsgStakeDelegate).Type():
		var msg types.MsgStakeDelegate
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgDelegate{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.DelegatorAddr, txMsg.ValidatorAddr)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.DelegatorAddress.String()
		docTx.To = msg.ValidatorAddress.String()
		docTx.Amount = []store.Coin{types.ParseCoin(msg.Amount.String())}
		docTx.Type = constant.TxTypeStakeDelegate

	case new(types.MsgStakeBeginUnbonding).Type():
		var msg types.MsgStakeBeginUnbonding
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgBeginUnbonding{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.DelegatorAddr, txMsg.ValidatorAddr)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.DelegatorAddress.String()
		docTx.To = msg.ValidatorAddress.String()
		docTx.Amount = []store.Coin{types.ParseCoin(msg.Amount.String())}
		docTx.Type = constant.TxTypeStakeBeginUnbonding
	case new(types.MsgBeginRedelegate).Type():
		var msg types.MsgBeginRedelegate
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgBeginRedelegate{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.DelegatorAddr, txMsg.ValidatorSrcAddr, txMsg.ValidatorDstAddr)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.ValidatorSrcAddress.String()
		docTx.To = msg.ValidatorDstAddress.String()
		docTx.Amount = []store.Coin{types.ParseCoin(msg.Amount.String())}
		docTx.Type = constant.TxTypeBeginRedelegate
	case new(types.MsgUnjail).Type():
		var msg types.MsgUnjail
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgUnjail{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.ValidatorAddr)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.ValidatorAddr.String()
		docTx.Type = constant.TxTypeUnjail
	default:
		ok = false
	}
	return docTx, ok
}