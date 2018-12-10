package handler

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
	"sync"
)

// Delegation represents the bond with tokens held by an account.  It is
// owned by one delegator, and is associated with the voting power of one
// pubKey.
type tempDelegation struct {
	Shares         float64
	OriginalShares string
	Height         int64 // Last height bond updated
}

// save Tx document into collection
func SaveTx(docTx document.CommonTx, mutex sync.Mutex) {
	var (
		methodName = "SaveTx"
	)
	logger.Debug("Start", logger.String("method", methodName))

	// save common docTx document
	saveCommonTx := func(commonTx document.CommonTx) {
		//save tx
		err := store.Save(commonTx)
		if err != nil {
			logger.Error("Save commonTx failed", logger.Any("Tx", commonTx), logger.String("err", err.Error()))
		}
		//save tx_msg
		msg := commonTx.Msg
		if msg != nil {
			txMsg := document.TxMsg{
				Hash:    docTx.TxHash,
				Type:    msg.Type(),
				Content: msg.String(),
			}
			store.Save(txMsg)
		}
		handleProposal(commonTx)
	}

	saveCommonTx(docTx)
	saveValidatorAndDelegator(docTx, mutex)
	logger.Debug("End", logger.String("method", methodName))
}

// save or update validator or delegator info
// by parse stake tx

// Different transaction types correspond to different operations TODO
//TxTypeStakeCreateValidator
//	1:insert new validator (---> CompareAndUpdateValidators)
//	2:insert delegator
//
//TxTypeStakeEditValidator
//	1:update validator
//
//TxTypeStakeDelegate
//	1:update validator (---> CompareAndUpdateValidators)
//	2:insert delegator(or update delegator existed )
//
//TxTypeStakeBeginUnbonding
//	1:update validator (---> CompareAndUpdateValidators)
//	2:update delegator
//
//TxTypeBeginRedelegate
//	1:update validator(src,dest) (---> CompareAndUpdateValidators)
//	2:update delegator
func saveValidatorAndDelegator(docTx document.CommonTx, mutex sync.Mutex) {
	var (
		methodName = "saveValidatorAndDelegator: "
		valAddress string
		delAddress string
		candidate  document.Candidate
		delegator  document.Delegator
	)

	logger.Debug("Start", logger.String("method", methodName))
	txType := GetTxType(docTx)
	if txType == "" {
		logger.Error("Tx invalid", logger.Any("Tx", docTx))
		return
	}

	switch txType {
	case constant.TxTypeStakeCreateValidator:
		valAddress = docTx.From
		delAddress = valAddress
	case constant.TxTypeStakeEditValidator:
		valAddress = docTx.From
		break
	case constant.TxTypeStakeDelegate, constant.TxTypeStakeBeginUnbonding,
		constant.TxTypeStakeCompleteUnbonding:
		valAddress = docTx.To
		delAddress = docTx.From
		break
	}

	if valAddress == "" {
		return
	}

	// get validator
	validator, err := helper.GetValidator(valAddress)

	if err != nil {
		logger.Error("validator not existed", logger.String("validator", valAddress))
		return
	}

	if validator.OperatorAddr == nil {
		// validator not exist
		candidate = document.Candidate{
			Address: valAddress,
		}
	} else {
		candidate = BuildValidatorDocument(validator)
	}

	// get delegator
	if delAddress != "" {

		// get delegation
		delegation, err := buildDelegation(delAddress, valAddress)
		if err != nil {
			logger.Error("get delegation failed by valAddr  and delAddr", logger.String("valAddress", valAddress), logger.String("delAddress", delAddress))
			return
		}

		// get unbondingDelegation
		ud, err := buildUnbondingDelegation(delAddress, valAddress)
		if err != nil {
			logger.Error("get unbonding delegation failed by valAddr  and delAddr", logger.String("valAddress", valAddress), logger.String("delAddress", delAddress))
			return
		}

		delegator = document.Delegator{
			Address:       delAddress,
			ValidatorAddr: valAddress,

			Shares:         delegation.Shares,
			OriginalShares: delegation.OriginalShares,
			BondedHeight:   delegation.Height,

			UnbondingDelegation: document.UnbondingDelegation{
				CreationHeight: ud.CreationHeight,
				MinTime:        ud.MinTime,
				InitialBalance: ud.InitialBalance,
				Balance:        ud.Balance,
			},
		}
	}

	mutex.Lock()
	defer mutex.Unlock()
	logger.Info("mutex.Lock()")

	// update or delete validator
	if candidate.PubKey == "" {
		store.Delete(candidate)
		logger.Info("delete candidate", logger.String("Address", candidate.Address))
	} else {
		store.SaveOrUpdate(candidate)
		logger.Info("saveOrUpdate candidate", logger.String("Address", candidate.Address))
	}

	// update or delete delegator
	if delAddress != "" {
		if delegator.BondedHeight < 0 &&
			delegator.UnbondingDelegation.CreationHeight < 0 {
			store.Delete(delegator)
			logger.Info("delete delegator", logger.String("Address", delegator.Address), logger.String("ValidatorAddr", delegator.ValidatorAddr))
		} else {
			store.SaveOrUpdate(delegator)
			logger.Info("saveOrUpdate delegator", logger.String("Address", delegator.Address), logger.String("ValidatorAddr", delegator.ValidatorAddr))
		}
	}
	logger.Debug("End", logger.String("method", methodName))
	logger.Info("release lock")
}

func buildDelegation(delAddress, valAddress string) (tempDelegation, error) {
	var (
		res tempDelegation
	)

	d, err := helper.GetDelegation(delAddress, valAddress)

	if err != nil {
		return res, err
	}

	if d.DelegatorAddr == nil {
		// represents delegation is nil
		res.Height = -1
		return res, nil
	}

	floatShares := helper.ParseFloat(d.Shares.String())
	res = tempDelegation{
		Shares:         floatShares,
		OriginalShares: d.Shares.String(),
		Height:         d.Height,
	}

	return res, nil
}

func buildUnbondingDelegation(delAddress, valAddress string) (
	document.UnbondingDelegation, error) {
	var (
		res document.UnbondingDelegation
	)

	ud, err := helper.GetUnbondingDelegation(delAddress, valAddress)

	if err != nil {
		return res, nil
	}
	// doesn't have unbonding delegation
	if ud.DelegatorAddr == nil {
		// represents unbonding delegation is nil
		res.CreationHeight = -1
		return res, nil
	}

	initBalance := types.BuildCoins(types.SdkCoins{ud.InitialBalance})
	balance := types.BuildCoins(types.SdkCoins{ud.Balance})
	res = document.UnbondingDelegation{
		CreationHeight: ud.CreationHeight,
		MinTime:        ud.MinTime.Unix(),
		InitialBalance: initBalance,
		Balance:        balance,
	}

	return res, nil
}
