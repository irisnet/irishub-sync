package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
	"github.com/tendermint/tendermint/types"
)

// compare validatorSet stored in tendermint and validatorSet stored in db
// if tmValidatorSet not qual dbValidatorSet, execute two step as follow.
// first, remove all validators stored in db
// second, store latest validators which query from sdk store into db
// note: this function isn't thread safe, should be invoked during watch block
//       not fast sync
func CompareAndUpdateValidators(tmVals []*types.Validator) {
	var (
		methodName     = "CompareAndUpdateValidators"
		tmValidatorSet []string
		dbValidatorSet []string

		candidateModel document.Candidate
		candidates     []document.Candidate

		kvs []sdk.KVPair
	)

	// get validatorSets from tendermint
	for _, v := range tmVals {
		tmValidatorSet = append(tmValidatorSet, v.Address.String())
	}

	// get unRevoke validatorSets from db
	dbVals, err := candidateModel.GetUnRevokeValidators()
	if err != nil {
		logger.Error.Printf("%v: err is %v\n", methodName, err)
	}
	for _, v := range dbVals {
		dbValidatorSet = append(dbValidatorSet, v.PubKeyAddr)
	}

	// tmValidatorSet not equal storeValidatorSet
	if !compareSlice(tmValidatorSet, dbValidatorSet) {
		logger.Info.Printf("%v: vlidatorSet changes, tmValSet is %v, dbValSet is %v\n",
			methodName, tmValidatorSet, dbValidatorSet)

		// remove all data which stored in db
		err := candidateModel.RemoveCandidates()
		if err != nil {
			logger.Error.Printf("%v: err is %v\n", methodName, err)
		}

		// store latest validator data
		// get latest validators through query sdk store
		keys := stake.ValidatorsKey
		resRaw, err := helper.Query(keys, constant.StoreNameStake, "subspace")

		if err != nil {
			logger.Error.Printf("%v: err is %v\n", methodName, err)
		}
		codec.Cdc.MustUnmarshalBinary(resRaw, &kvs)
		for _, v := range kvs {
			var (
				validator stake.Validator
			)
			err2 := codec.Cdc.UnmarshalBinary(v.Value, &validator)
			if err2 != nil {
				logger.Error.Printf("%v: err is %v\n", methodName, err2)
			}

			// build validator document struct by stake.validator
			doc := BuildValidatorDocument(validator)
			candidates = append(candidates, doc)
		}

		// store latest validators into db
		err3 := candidateModel.SaveAll(candidates)
		if err3 != nil {
			logger.Error.Printf("%v: err is %v\n", methodName, err3)
		}
	} else {
		logger.Info.Printf("%v: validatorSet not change\n", methodName)
	}
}

func compareSlice(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range b {
		if !sliceContains(a, b[i]) {
			return false
		}
	}

	return true
}

// contains method for a slice
func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func BuildValidatorDocument(v stake.Validator) document.Candidate {
	description := document.Description{
		Moniker:  v.Description.Moniker,
		Identity: v.Description.Identity,
		Website:  v.Description.Website,
		Details:  v.Description.Details,
	}

	floatShares, _ := v.PoolShares.Amount.Float64()
	doc := document.Candidate{
		Address:        v.Owner.String(),
		PubKey:         helper.BuildHex(v.PubKey.Bytes()),
		PubKeyAddr:     v.PubKey.Address().String(),
		Revoked:        v.Revoked,
		Shares:         floatShares,
		OriginalShares: v.PoolShares.Amount.String(),
		Description:    description,
		BondHeight:     v.BondHeight,
	}

	doc.VotingPower = doc.Shares

	return doc
}
