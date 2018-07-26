package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
	"sync"
)

// save Tx document into collection
func SaveTx(docTx store.Docs, mutex sync.Mutex) {
	var (
		methodName = "SaveTx: "
		valAddress string
		delAddress string
	)

	// save docTx document into database
	storeTxDocFunc := func(doc store.Docs) {
		err := store.Save(doc)
		if err != nil {
			logger.Error.Printf("%v Save failed. doc is %+v, err is %v",
				methodName, doc, err.Error())
		}
	}

	// save common docTx document
	saveCommonTx := func(commonTx document.CommonTx) {
		err := store.Save(commonTx)
		if err != nil {
			logger.Error.Printf("%v Save commonTx failed. doc is %+v, err is %v",
				methodName, commonTx, err.Error())
		}
	}

	txType := GetTxType(docTx)
	if txType == "" {
		logger.Error.Printf("%v get docTx type failed, docTx is %v\n",
			methodName, docTx)
		return
	}

	saveCommonTx(buildCommonTxData(docTx, txType))

	switch txType {
	case constant.TxTypeBank:
		break
	case constant.TxTypeStakeCreate:
		docTx, r := docTx.(document.StakeTxDeclareCandidacy)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		storeTxDocFunc(docTx)

		valAddress = docTx.ValidatorAddr
		break
	case constant.TxTypeStakeEdit:
		docTx, r := docTx.(document.StakeTxEditCandidacy)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		storeTxDocFunc(docTx)

		valAddress = docTx.ValidatorAddr
		break
	case constant.TxTypeStakeDelegate:
		docTx, r := docTx.(document.StakeTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		storeTxDocFunc(docTx)

		valAddress = docTx.ValidatorAddr
		delAddress = docTx.DelegatorAddr
		break
	case constant.TxTypeStakeUnbond:
		docTx, r := docTx.(document.StakeTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		storeTxDocFunc(docTx)

		valAddress = docTx.ValidatorAddr
		delAddress = docTx.DelegatorAddr
		break
	}

	if valAddress != "" {
		var (
			candidate document.Candidate
			delegator document.Delegator
		)

		// prepare validator data
		validator, err := getValidator(valAddress)

		if err != nil {
			logger.Error.Printf("%v get validator failed by valAddr %v\n", methodName, valAddress)
			return
		}

		if validator.Owner == nil {
			// validator not exist
			candidate = document.Candidate{
				Address: valAddress,
			}
		} else {
			candidate = BuildValidatorDocument(validator)

			// prepare delegator data
			if delAddress != "" {
				delegation, err := getDelegation(delAddress, valAddress)

				if err != nil {
					logger.Error.Printf("%v get delegation failed by valAddr %v and delAddr %v\n", methodName, valAddress, delAddress)
					return
				}

				if delegation.DelegatorAddr == nil {
					// can't get delegation when delegator unbond all token
					delegator = document.Delegator{
						Address:        delAddress,
						ValidatorAddr:  valAddress,
						Shares:         float64(0),
						OriginalShares: "",
					}
				} else {
					// delegation exist
					floatShares, _ := delegation.Shares.Float64()
					delegator = document.Delegator{
						Address:        delegation.DelegatorAddr.String(),
						ValidatorAddr:  delegation.ValidatorAddr.String(),
						Shares:         floatShares,
						OriginalShares: delegation.Shares.String(),
					}
				}
			}
		}

		mutex.Lock()
		logger.Info.Printf("%v saveOrUpdate cndidates get lock\n", methodName)

		// update or delete validator
		if candidate.PubKey == "" {
			store.Delete(candidate)
		} else {
			store.SaveOrUpdate(candidate)
		}

		// update or delete delegator
		if delAddress != "" {
			if delegator.OriginalShares == "" {
				store.Delete(delegator)
			} else {
				store.SaveOrUpdate(delegator)
			}
		}

		mutex.Unlock()
		logger.Info.Printf("%v saveOrUpdate cndidates release lock\n", methodName)
	}
}

// build common tx data through parse tx
func buildCommonTxData(docTx store.Docs, txType string) document.CommonTx {
	var commonTx document.CommonTx

	if txType == "" {
		txType = GetTxType(docTx)
	}
	switch txType {
	case constant.TxTypeBank:
		doc := docTx.(document.CommonTx)
		commonTx = document.CommonTx{
			TxHash: doc.TxHash,
			Time:   doc.Time,
			Height: doc.Height,
			From:   doc.From,
			To:     doc.To,
			Amount: doc.Amount,
			Type:   doc.Type,
			Fee:    doc.Fee,
			Status: doc.Status,
		}
		break
	case constant.TxTypeStakeCreate:
		doc := docTx.(document.StakeTxDeclareCandidacy)
		commonTx = document.CommonTx{
			TxHash: doc.TxHash,
			Time:   doc.Time,
			Height: doc.Height,
			From:   doc.ValidatorAddr,
			To:     "",
			Amount: []store.Coin{doc.Amount},
			Type:   doc.Type,
			Fee:    doc.Fee,
			Status: doc.Status,
		}
		break
	case constant.TxTypeStakeEdit:
		doc := docTx.(document.StakeTxEditCandidacy)
		commonTx = document.CommonTx{
			TxHash: doc.TxHash,
			Time:   doc.Time,
			Height: doc.Height,
			From:   doc.ValidatorAddr,
			To:     "",
			Amount: []store.Coin{doc.Amount},
			Type:   doc.Type,
			Fee:    doc.Fee,
			Status: doc.Status,
		}
		break
	case constant.TxTypeStakeDelegate, constant.TxTypeStakeUnbond:
		doc := docTx.(document.StakeTx)
		commonTx = document.CommonTx{
			TxHash: doc.TxHash,
			Time:   doc.Time,
			Height: doc.Height,
			From:   doc.DelegatorAddr,
			To:     doc.ValidatorAddr,
			Amount: []store.Coin{doc.Amount},
			Type:   doc.Type,
			Fee:    doc.Fee,
			Status: doc.Status,
		}
		break
	}

	return commonTx
}

// get validator
func getValidator(valAddr string) (stake.Validator, error) {
	var (
		validatorAddr sdk.Address
		err           error
		res           stake.Validator
	)

	validatorAddr, err = sdk.GetAccAddressHex(valAddr)

	resRaw, err := helper.Query(stake.GetValidatorKey(validatorAddr), constant.StoreNameStake, constant.StoreDefaultEndPath)
	if err != nil || resRaw == nil {
		return res, err
	}

	codec.Cdc.MustUnmarshalBinary(resRaw, &res)

	return res, err
}

// get delegation
func getDelegation(delAddr, valAddr string) (stake.Delegation, error) {
	var (
		delegatorAddr sdk.Address
		validatorAddr sdk.Address
		err           error

		res stake.Delegation
	)

	delegatorAddr, err = sdk.GetAccAddressHex(delAddr)

	if err != nil {
		return res, err
	}

	validatorAddr, err = sdk.GetAccAddressHex(valAddr)

	if err != nil {
		return res, err
	}
	cdc := codec.Cdc
	key := stake.GetDelegationKey(delegatorAddr, validatorAddr, cdc)

	resRaw, err := helper.Query(key, constant.StoreNameStake, constant.StoreDefaultEndPath)

	if err != nil || resRaw == nil {
		return res, err
	}

	err = cdc.UnmarshalBinary(resRaw, &res)

	if err != nil {
		return res, err
	}

	return res, err
}
