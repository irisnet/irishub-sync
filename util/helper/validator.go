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

	if err != nil {
		logger.Error("helper.GetValidators err ", logger.String("err", err.Error()))
		return
	}

	cdc.MustUnmarshalBinaryLengthPrefixed(resRaw, &kvs)

	for _, v := range kvs {
		addr := v.Key[1:]
		validator, err2 := types.UnmarshalValidator(cdc, addr, v.Value)

		if err2 != nil {
			logger.Error("types.UnmarshalValidator", logger.String("err", err2.Error()))
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

// get delegation
func GetDelegation(delAddr, valAddr string) (types.Delegation, error) {
	var (
		delegatorAddr types.AccAddress
		validatorAddr types.ValAddress
		err           error

		res types.Delegation
	)
	cdc := types.GetCodec()

	delegatorAddr, err = types.AccAddressFromBech32(delAddr)

	if err != nil {
		return res, err
	}

	validatorAddr, err = types.ValAddressFromBech32(valAddr)

	if err != nil {
		return res, err
	}
	key := types.GetDelegationKey(delegatorAddr, validatorAddr) //TODO

	resRaw, err := Query(key, constant.StoreNameStake, constant.StoreDefaultEndPath)

	if err != nil || resRaw == nil {
		return res, err
	}

	res, err = types.UnmarshalDelegation(cdc, key, resRaw)

	if err != nil {
		return res, err
	}

	return res, err
}

// get unbonding delegation
func GetUnbondingDelegation(delAddr, valAddr string) (types.UnbondingDelegation, error) {
	var (
		delegatorAddr types.AccAddress
		validatorAddr types.ValAddress
		err           error

		res types.UnbondingDelegation
	)

	cdc := types.GetCodec()

	delegatorAddr, err = types.AccAddressFromBech32(delAddr)

	if err != nil {
		return res, err
	}

	validatorAddr, err = types.ValAddressFromBech32(valAddr)

	if err != nil {
		return res, err
	}

	key := types.GetUBDKey(delegatorAddr, validatorAddr) //TODO ValAddressFromBech32

	resRaw, err := Query(key, constant.StoreNameStake, constant.StoreDefaultEndPath)

	if err != nil || resRaw == nil {
		return res, err
	}

	res = types.MustUnmarshalUBD(cdc, key, resRaw)

	return res, nil
}
