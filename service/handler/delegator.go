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

// save or update validator or delegator info
// by parse stake tx

// Different transaction types correspond to different operations
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
//	2:update delegator(src,dest)
func SaveOrUpdateDelegator(docTx document.CommonTx, mutex sync.Mutex) {

	logger.Debug("Start", logger.String("method", "saveDelegator"))

	switch docTx.Type {
	case constant.TxTypeStakeCreateValidator:
		modifyDelegator(docTx.From, docTx.From)
		break
	case constant.TxTypeStakeEditValidator:
		updateValidator(docTx.From)
		break
	case constant.TxTypeStakeDelegate, constant.TxTypeStakeBeginUnbonding:
		modifyDelegator(docTx.From, docTx.To)
		break
	case constant.TxTypeBeginRedelegate:
		delAddress := docTx.From
		msg := docTx.Msg.(types.BeginRedelegate)
		valSrcAddr := msg.ValidatorSrcAddr
		valDstAddr := msg.ValidatorDstAddr

		modifyDelegator(delAddress, valSrcAddr)
		modifyDelegator(delAddress, valDstAddr)
		break
	}

	logger.Debug("End", logger.String("method", "saveDelegator"))
}

func modifyDelegator(delAddress, valAddress string) {
	// get delegation
	delegation, err := buildDelegation(delAddress, valAddress)
	if err != nil {
		logger.Error("get delegation failed", logger.String("valAddress", valAddress), logger.String("delAddress", delAddress))
		return
	}

	// get unbondingDelegation
	ud, err := buildUnbondingDelegation(delAddress, valAddress)
	if err != nil {
		logger.Error("get unbonding delegation", logger.String("valAddress", valAddress), logger.String("delAddress", delAddress))
		return
	}

	delegator := document.Delegator{
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

	if delegator.BondedHeight < 0 &&
		delegator.UnbondingDelegation.CreationHeight < 0 {
		store.Delete(delegator)
		logger.Info("delete delegator", logger.String("Address", delegator.Address), logger.String("ValidatorAddr", delegator.ValidatorAddr))
	} else {
		store.SaveOrUpdate(delegator)
		logger.Info("saveOrUpdate delegator", logger.String("Address", delegator.Address), logger.String("ValidatorAddr", delegator.ValidatorAddr))
	}
}

func updateValidator(valAddress string) {
	//var canCollection  document.Candidate

	validator, err := helper.GetValidator(valAddress)
	if err != nil {
		logger.Error("validator not existed", logger.String("validator", valAddress))
		return
	}

	editValidator := BuildValidatorDocument(validator)
	//candidate := canCollection.GetValidator(valAddress)
	//editValidator.Rank = candidate.Rank
	if err := store.Update(editValidator); err != nil {
		logger.Error("update candidate error", logger.String("address", valAddress))
	}
	logger.Info("Update candidate success", logger.String("Address", valAddress))
}

func buildDelegation(delAddress, valAddress string) (res tempDelegation, err error) {
	d := helper.GetDelegation(delAddress, valAddress)

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

	ud := helper.GetUnbondingDelegation(delAddress, valAddress)

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

// Delegation represents the bond with tokens held by an account.  It is
// owned by one delegator, and is associated with the voting power of one
// pubKey.
type tempDelegation struct {
	Shares         float64
	OriginalShares string
	Height         int64 // Last height bond updated
}
