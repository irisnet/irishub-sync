package handler

import (
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/msg/guardian"
	"github.com/irisnet/irishub-sync/msg/distribution"
	"github.com/irisnet/irishub-sync/util/helper"
	"encoding/json"
	"github.com/irisnet/irishub-sync/msg/iservice"
)

//when new address not found or the tx is not success,this address will  not be collected
func saveNewAccount(tx *document.CommonTx) {
	var accountModel document.Account
	if document.TxStatusSuccess != tx.Status {
		return
	}
	switch tx.Type {
	case constant.TxTypeTransfer:
		accountModel.Address = tx.To
	case constant.TxTypeAddTrustee:
		if len(tx.Msgs) > 0 {
			msgData := guardian.DocTxMsgAddTrustee{}
			if err := json.Unmarshal([]byte(helper.ToJson(tx.Msgs[0].Msg)), &msgData); err == nil {
				accountModel.Address = msgData.Address
			}
		}

	case constant.TxTypeSetWithdrawAddress:
		if len(tx.Msgs) > 0 {
			msgData := distribution.DocTxMsgSetWithdrawAddress{}
			if err := json.Unmarshal([]byte(helper.ToJson(tx.Msgs[0].Msg)), &msgData); err == nil {
				accountModel.Address = msgData.WithdrawAddr
			}
		}
	case constant.TxTypeServiceSetWithdrawAddress:
		if len(tx.Msgs) > 0 {
			msgData := iservice.DocMsgServiceSetWithdrawAddress{}
			if err := json.Unmarshal([]byte(helper.ToJson(tx.Msgs[0].Msg)), &msgData); err == nil {
				accountModel.Address = msgData.WithdrawAddress
			}
		}
	}
	if accountModel.Address == "" {
		return
	}
	if err := accountModel.SaveAddress(accountModel.Address); err != nil {
		logger.Warn("Save new account address failed", logger.String("err", err.Error()))
	}
}
