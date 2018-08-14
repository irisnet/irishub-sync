package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cosmos/cosmos-sdk/x/stake/types"
	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
	"sync"
)

// save Tx document into collection
func SaveTx(docTx document.CommonTx, mutex sync.Mutex) {
	var (
		methodName = "SaveTx: "
	)
	logger.Info.Printf("Start %v\n", methodName)

	// save common docTx document
	saveCommonTx := func(commonTx document.CommonTx) {
		err := store.Save(commonTx)
		if err != nil {
			logger.Error.Printf("%v Save commonTx failed. doc is %+v, err is %v",
				methodName, commonTx, err.Error())
		}
	}

	saveCommonTx(docTx)

	saveValidatorAndDelegator(docTx, mutex)

	logger.Info.Printf("End %v\n", methodName)
}

// save or update validator or delegator info
// by parse stake tx
func saveValidatorAndDelegator(docTx document.CommonTx, mutex sync.Mutex) {
	var (
		methodName = "saveValidatorAndDelegator: "
		valAddress string
		delAddress string
		candidate  document.Candidate
		delegator  document.Delegator
	)

	txType := GetTxType(docTx)
	if txType == "" {
		logger.Error.Printf("%v get docTx type failed, docTx is %v\n",
			methodName, docTx)
		return
	}

	switch txType {
	case constant.TxTypeStakeCreateValidator, constant.TxTypeStakeEditValidator:
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
	validator, err := getValidator(valAddress)

	if err != nil {
		logger.Error.Printf("%v: get validator failed by valAddr %v\n", methodName, valAddress)
		return
	}

	if validator.Owner == nil {
		// validator not exist
		candidate = document.Candidate{
			Address: valAddress,
		}
	} else {
		candidate = BuildValidatorDocument(validator)
	}

	// get delegator
	if delAddress != "" {
		delegation, err := getDelegation(delAddress, valAddress)

		if err != nil {
			logger.Error.Printf("%v: get delegation failed by valAddr %v and delAddr %v\n", methodName, valAddress, delAddress)
			return
		}

		if delegation.DelegatorAddr == nil {
			logger.Info.Printf("%v: delegation is nil\n", methodName)
			// can't get delegation when delegator unbond all token
			delegator = document.Delegator{
				Address:        delAddress,
				ValidatorAddr:  valAddress,
				Shares:         float64(-1),
				OriginalShares: "",
			}
		} else {
			// delegation exist
			floatShares, _ := delegation.Shares.Float64()
			delegator = document.Delegator{
				Address:        delegation.DelegatorAddr.String(),
				ValidatorAddr:  delegation.ValidatorAddr.String(),
				Shares:         floatShares,
				OriginalShares: delegation.Shares.RatString(),
				Height:         delegation.Height,
			}
		}
	}

	mutex.Lock()
	logger.Info.Printf("%v saveOrUpdate vals and dels get lock\n", methodName)

	// update or delete validator
	if candidate.PubKey == "" {
		store.Delete(candidate)
		logger.Info.Printf("%v delete candidate, addr is %v\n", methodName, candidate.Address)
	} else {
		store.SaveOrUpdate(candidate)
		logger.Info.Printf("%v saveOrUpdate candidate, addr is %v\n", methodName, candidate.Address)
	}

	// update or delete delegator
	if delAddress != "" {
		if delegator.Shares <= float64(0) {
			store.Delete(delegator)
			logger.Info.Printf("%v delete delegator, delVar is %v, valAddr is %v\n",
				methodName, delegator.Address, delegator.ValidatorAddr)
		} else {
			store.SaveOrUpdate(delegator)
			logger.Info.Printf("%v saveOrUpdate delegator, delVar is %v, valAddr is %v\n",
				methodName, delegator.Address, delegator.ValidatorAddr)
		}
	}

	mutex.Unlock()
	logger.Info.Printf("%v saveOrUpdate vals and dels release lock\n", methodName)
}

// get validator
func getValidator(valAddr string) (stake.Validator, error) {
	var (
		validatorAddr sdk.AccAddress
		err           error
		res           stake.Validator
	)

	validatorAddr, err = sdk.AccAddressFromBech32(valAddr)

	resRaw, err := helper.Query(stake.GetValidatorKey(validatorAddr), constant.StoreNameStake, constant.StoreDefaultEndPath)
	if err != nil || resRaw == nil {
		return res, err
	}

	res = types.MustUnmarshalValidator(codec.Cdc, validatorAddr, resRaw)

	return res, err
}

// get delegation
func getDelegation(delAddr, valAddr string) (stake.Delegation, error) {
	var (
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.AccAddress
		err           error

		res stake.Delegation
	)

	delegatorAddr, err = sdk.AccAddressFromBech32(delAddr)

	if err != nil {
		return res, err
	}

	validatorAddr, err = sdk.AccAddressFromBech32(valAddr)

	if err != nil {
		return res, err
	}
	cdc := codec.Cdc
	key := stake.GetDelegationKey(delegatorAddr, validatorAddr)

	resRaw, err := helper.Query(key, constant.StoreNameStake, constant.StoreDefaultEndPath)

	if err != nil || resRaw == nil {
		return res, err
	}

	res, err = types.UnmarshalDelegation(cdc, key, resRaw)

	if err != nil {
		return res, err
	}

	return res, err
}
