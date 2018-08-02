package handler

import (
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
	"sync"
	"time"
)

// save account
func SaveAccount(docTx store.Docs, mutex sync.Mutex) {
	var (
		address    string
		updateTime time.Time
		height     int64
		methodName = "SaveAccount: "
	)
	logger.Info.Printf("Start %v\n", methodName)

	// save account
	fun := func(address string, updateTime time.Time, height int64) {
		account := document.Account{
			Address: address,
			Time:    updateTime,
			Height:  height,
		}

		err := store.Save(account)

		if err != nil && err.Error() != "Record exists" {
			logger.Error.Printf("%v Record exists, account is %v, err is %s\n",
				methodName, account.Address, err.Error())
		}
	}

	txType := GetTxType(docTx)
	if txType == "" {
		logger.Error.Printf("%v get docTx type failed, docTx is %v\n",
			methodName, docTx)
		return
	}

	switch txType {
	case constant.TxTypeTransfer:
		docTx, r := docTx.(document.CommonTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		updateTime = docTx.Time
		height = docTx.Height

		fun(docTx.From, updateTime, height)
		fun(docTx.To, updateTime, height)
		break
	case constant.TxTypeStakeCreateValidator:
		docTx, r := docTx.(document.StakeTxDeclareCandidacy)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		address = docTx.ValidatorAddr
		updateTime = docTx.Time
		height = docTx.Height

		fun(address, updateTime, height)
		if docTx.DelegatorAddr != "" {
			fun(docTx.DelegatorAddr, updateTime, height)
		}
		break
	case constant.TxTypeStakeEditValidator:
		docTx, r := docTx.(document.StakeTxEditCandidacy)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		address = docTx.ValidatorAddr
		updateTime = docTx.Time
		height = docTx.Height

		fun(address, updateTime, height)
		break
	case constant.TxTypeStakeDelegate, constant.TxTypeStakeBeginUnbonding,
		constant.TxTypeStakeCompleteUnbonding:
		stakeTx, r := docTx.(document.StakeTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		updateTime = stakeTx.Time
		height = stakeTx.Height

		fun(stakeTx.ValidatorAddr, updateTime, height)
		fun(stakeTx.DelegatorAddr, updateTime, height)
		break
	}

	logger.Info.Printf("End %v\n", methodName)
}

// update account balance
func UpdateBalance(docTx store.Docs, mutex sync.Mutex) {
	var (
		methodName = "UpdateBalance: "
	)
	logger.Info.Printf("Start %v\n", methodName)

	fun := func(address string) {
		account, err := document.QueryAccount(address)
		if err != nil {
			logger.Error.Printf("%v updateAccountBalance failed, account is %v and err is %v",
				methodName, account, err.Error())
			return
		}

		// query balance of account
		account.Amount = helper.QueryAccountBalance(address)
		if err := store.Update(account); err != nil {
			logger.Error.Printf("%v account:[%q] balance update failed,%s\n",
				methodName, account.Address, err)
		}
	}

	txType := GetTxType(docTx)
	if txType == "" {
		logger.Error.Printf("%v get docTx type failed, docTx is %v\n",
			methodName, docTx)
		return
	}

	switch txType {
	case constant.TxTypeTransfer:
		docTx, r := docTx.(document.CommonTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		fun(docTx.From)
		fun(docTx.To)
		break
	case constant.TxTypeStakeCreateValidator:
		docTx, r := docTx.(document.StakeTxDeclareCandidacy)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		fun(docTx.ValidatorAddr)
		if docTx.DelegatorAddr != "" {
			fun(docTx.DelegatorAddr)
		}
		break
	case constant.TxTypeStakeEditValidator:
		docTx, r := docTx.(document.StakeTxEditCandidacy)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		fun(docTx.ValidatorAddr)
		break
	case constant.TxTypeStakeDelegate, constant.TxTypeStakeBeginUnbonding,
		constant.TxTypeStakeCompleteUnbonding:
		docTx, r := docTx.(document.StakeTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		fun(docTx.ValidatorAddr)
		fun(docTx.DelegatorAddr)
		break
	}

	logger.Info.Printf("End %v\n", methodName)
}
