package helper

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

func GetValidators() (validators []types.StakeValidator) {
	keys := types.ValidatorsKey
	cdc := types.GetCodec()
	var kvs []types.KVPair

	resRaw, err := Query(keys, constant.StoreNameStake, "subspace")

	if err != nil || len(resRaw) == 0 {
		logger.Error("GetValidators Failed ", logger.String("err", err.Error()))
		return
	}

	err = cdc.UnmarshalBinaryLengthPrefixed(resRaw, &kvs)
	if err != nil {
		logger.Error("UnmarshalBinaryLengthPrefixed validators err ", logger.String("err", err.Error()))
		return
	}

	for _, v := range kvs {
		addr := v.Key[1:]
		validator, err2 := types.UnmarshalValidator(cdc, addr, v.Value)

		if err2 != nil {
			logger.Error("types.UnmarshalValidator", logger.String("err", err2.Error()))
			continue
		}

		validators = append(validators, validator)
	}
	return validators
}

// get validator
func GetValidator(valAddr string) (types.StakeValidator, error) {
	var (
		validatorAddr types.ValAddress
		err           error
		res           types.StakeValidator
	)

	cdc := types.GetCodec()

	validatorAddr, err = types.ValAddressFromBech32(valAddr)

	resRaw, err := Query(types.GetValidatorKey(validatorAddr), constant.StoreNameStake, constant.StoreDefaultEndPath) //TODO
	if err != nil || resRaw == nil {
		return res, err
	}

	res = types.MustUnmarshalValidator(cdc, validatorAddr, resRaw)

	return res, err
}

// Query a delegation based on address and validator address
func GetDelegation(delAddr, valAddr string) (res types.Delegation) {
	var (
		validatorAddr types.ValAddress
		err           error
	)
	cdc := types.GetCodec()

	delegatorAddr, _ := types.AccAddressFromBech32(delAddr)
	validatorAddr, _ = types.ValAddressFromBech32(valAddr)

	key := types.GetDelegationKey(delegatorAddr, validatorAddr)

	resRaw, err := Query(key, constant.StoreNameStake, constant.StoreDefaultEndPath)

	if err != nil {
		logger.Error("helper.GetDelegation err ", logger.String("delAddr", delAddr))
		return
	} else if resRaw == nil {
		logger.Info("delegator don't exist delegation on validator", logger.String("delAddr", delAddr), logger.String("valAddr", valAddr))
		return
	}

	res = types.MustUnmarshalDelegation(cdc, key, resRaw)
	return res
}

//Query all delegations made from one delegator
func GetDelegations(delAddr string) (delegations []types.Delegation) {

	delegatorAddr, err := types.AccAddressFromBech32(delAddr)
	key := types.GetDelegationsKey(delegatorAddr)
	resKVs, err := QuerySubspace(key, constant.StoreNameStake)

	if err != nil {
		logger.Error("helper.GetDelegations err ", logger.String("delAddr", delAddr))
		return
	} else if resKVs == nil {
		logger.Info("delegator don't exist delegation", logger.String("delAddr", delAddr))
		return
	}

	cdc := types.GetCodec()

	for _, kv := range resKVs {
		delegation := types.MustUnmarshalDelegation(cdc, kv.Key, kv.Value)
		delegations = append(delegations, delegation)
	}
	return
}

// GetCmdQueryUnbondingDelegation implements the command to query a single unbonding-delegation record.
func GetUnbondingDelegation(delAddr, valAddr string) (res types.UnbondingDelegation) {
	cdc := types.GetCodec()

	delegatorAddr, _ := types.AccAddressFromBech32(delAddr)
	validatorAddr, _ := types.ValAddressFromBech32(valAddr)

	key := types.GetUBDKey(delegatorAddr, validatorAddr)

	resRaw, err := Query(key, constant.StoreNameStake, constant.StoreDefaultEndPath)

	if err != nil {
		logger.Error("helper.GetDelegations err ", logger.String("delAddr", delAddr))
		return
	} else if resRaw == nil {
		logger.Info("delegator don't exist unbondingDelegation", logger.String("delAddr", delAddr), logger.String("valAddr", valAddr))
		return
	}

	res = types.MustUnmarshalUBD(cdc, key, resRaw)

	return res
}

//Query all unbonding-delegations records for one delegator
func GetUnbondingDelegations(delAddr string) (ubds []types.UnbondingDelegation) {
	delegatorAddr, _ := types.AccAddressFromBech32(delAddr)

	cdc := types.GetCodec()
	key := types.GetUBDsKey(delegatorAddr)

	resKVs, err := QuerySubspace(key, constant.StoreNameStake)
	if err != nil {
		logger.Error("helper.GetDelegations err ", logger.String("delAddr", delAddr))
		return
	} else if resKVs == nil {
		logger.Info("delegator don't exist unbondingDelegation", logger.String("delAddr", delAddr))
		return
	}
	for _, kv := range resKVs {
		ubd := types.MustUnmarshalUBD(cdc, kv.Key, kv.Value)
		ubds = append(ubds, ubd)
	}
	return
}
