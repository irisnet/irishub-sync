package distribution

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
	case new(types.MsgSetWithdrawAddress).Type():
		var msg types.MsgSetWithdrawAddress
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgSetWithdrawAddress{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.DelegatorAddr, txMsg.WithdrawAddr)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.DelegatorAddress.String()
		docTx.To = msg.WithdrawAddress.String()
		docTx.Type = constant.TxTypeSetWithdrawAddress
	case new(types.MsgWithdrawDelegatorReward).Type():
		var msg types.MsgWithdrawDelegatorReward
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgWithdrawDelegatorReward{}
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
		docTx.Type = constant.TxTypeWithdrawDelegatorReward

	case new(types.MsgFundCommunityPool).Type():
		var msg types.MsgFundCommunityPool
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgFundCommunityPool{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTx.Msgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		docTx.Addrs = append(docTx.Addrs, txMsg.Depositor)
		docTx.Types = append(docTx.Types, txMsg.Type())
		if len(docTx.Msgs) > 1 {
			return docTx, true
		}
		docTx.From = msg.Depositor.String()
		docTx.Amount = types.ParseCoins(msg.Amount.String())
		docTx.Type = constant.TxTypeMsgFundCommunityPool
	case new(types.MsgWithdrawValidatorCommission).Type():
		var msg types.MsgWithdrawValidatorCommission
		data, _ := json.Marshal(msgData)
		json.Unmarshal(data, &msg)

		txMsg := DocTxMsgWithdrawValidatorCommission{}
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
		docTx.Type = constant.TxTypeMsgWithdrawValidatorCommission
	default:
		ok = false
	}
	return docTx, ok
}
