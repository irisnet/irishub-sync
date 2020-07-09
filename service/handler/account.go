package handler

import (
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/types/msg"
	"github.com/irisnet/irishub-sync/util/helper"
	"encoding/json"
)

func saveNewAccount(tx document.CommonTx) {
	var accountModel document.Account
	switch tx.Type {
	case constant.TxTypeTransfer:
		accountModel.Address = tx.To
	case constant.TxTypeAddTrustee:
		if len(tx.Msgs) > 0 {
			msgData := msg.DocTxMsgAddTrustee{}
			if err := json.Unmarshal([]byte(helper.ToJson(tx.Msgs[0].Msg)), &msgData); err == nil {
				accountModel.Address = msgData.Address
			}
		}

	case constant.TxTypeSetWithdrawAddress:
		if len(tx.Msgs) > 0 {
			msgData := msg.DocTxMsgSetWithdrawAddress{}
			if err := json.Unmarshal([]byte(helper.ToJson(tx.Msgs[0].Msg)), &msgData); err == nil {
				accountModel.Address = msgData.WithdrawAddr
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
