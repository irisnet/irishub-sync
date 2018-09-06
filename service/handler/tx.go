package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/module/logger"
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
		methodName = "SaveTx: "
	)
	logger.Info.Printf("Start %v\n", methodName)

	// save common docTx document
	saveCommonTx := func(commonTx document.CommonTx) {
		//save tx
		err := store.Save(commonTx)
		if err != nil {
			logger.Error.Printf("%v Save commonTx failed. doc is %+v, err is %v",
				methodName, commonTx, err.Error())
		}
		//save tx_msg
		msg := commonTx.Msg
		if msg != nil {
			txMsg := document.TxMsg{
				Hash:    docTx.TxHash,
				Content: msg.String(),
			}
			store.Save(txMsg)
		}
		handleProposal(commonTx)
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

		// get delegation
		delegation, err := buildDelegation(delAddress, valAddress)
		if err != nil {
			logger.Error.Printf("%v: get delegation failed by valAddr %v and delAddr %v\n",
				methodName, valAddress, delAddress)
			return
		}

		// get unbondingDelegation
		ud, err := buildUnbondingDelegation(delAddress, valAddress)
		if err != nil {
			logger.Error.Printf("%v: get unbonding delegation failed by valAddr %v and delAddr %v\n",
				methodName, valAddress, delAddress)
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
	logger.Info.Printf("%v get lock\n", methodName)

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
		if delegator.BondedHeight < 0 &&
			delegator.UnbondingDelegation.CreationHeight < 0 {
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
	logger.Info.Printf("%v release lock\n", methodName)
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

	floatShares, _ := d.Shares.Float64()
	res = tempDelegation{
		Shares:         floatShares,
		OriginalShares: d.Shares.RatString(),
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

	initBalance := types.BuildCoins(sdk.Coins{ud.InitialBalance})
	balance := types.BuildCoins(sdk.Coins{ud.Balance})
	res = document.UnbondingDelegation{
		CreationHeight: ud.CreationHeight,
		MinTime:        ud.MinTime,
		InitialBalance: initBalance,
		Balance:        balance,
	}

	return res, nil
}
