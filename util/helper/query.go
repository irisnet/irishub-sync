package helper

import (
	"fmt"
	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/pkg/errors"
)

// get validator
func GetValidator(valAddr string) (types.StakeValidator, error) {
	var (
		validatorAddr types.AccAddress
		err           error
		res           types.StakeValidator
	)

	validatorAddr, err = types.AccAddressFromBech32(valAddr)

	resRaw, err := Query(types.GetValidatorKey(validatorAddr), constant.StoreNameStake, constant.StoreDefaultEndPath)
	if err != nil || resRaw == nil {
		return res, err
	}

	res = types.MustUnmarshalValidator(codec.Cdc, validatorAddr, resRaw)

	return res, err
}

// get delegation
func GetDelegation(delAddr, valAddr string) (types.Delegation, error) {
	var (
		delegatorAddr types.AccAddress
		validatorAddr types.AccAddress
		err           error

		res types.Delegation
	)

	delegatorAddr, err = types.AccAddressFromBech32(delAddr)

	if err != nil {
		return res, err
	}

	validatorAddr, err = types.AccAddressFromBech32(valAddr)

	if err != nil {
		return res, err
	}
	cdc := codec.Cdc
	key := types.GetDelegationKey(delegatorAddr, validatorAddr)

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
		validatorAddr types.AccAddress
		err           error

		res types.UnbondingDelegation
	)

	delegatorAddr, err = types.AccAddressFromBech32(delAddr)

	if err != nil {
		return res, err
	}

	validatorAddr, err = types.AccAddressFromBech32(valAddr)

	if err != nil {
		return res, err
	}

	cdc := codec.Cdc
	key := types.GetUBDKey(delegatorAddr, validatorAddr)

	resRaw, err := Query(key, constant.StoreNameStake, constant.StoreDefaultEndPath)

	if err != nil || resRaw == nil {
		return res, err
	}

	res = types.MustUnmarshalUBD(cdc, key, resRaw)

	return res, nil
}

// Query from Tendermint with the provided storename and path
func Query(key types.HexBytes, storeName string, endPath string) (res []byte, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	client := GetClient()
	defer client.Release()

	rpcClient := client.Client
	if err != nil {
		return res, err
	}

	opts := types.ABCIQueryOptions{
		Height:  0,
		Trusted: true,
	}
	result, err := rpcClient.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}
	resp := result.Response
	if resp.Code != uint32(0) {
		return res, errors.Errorf("Query failed: (%d) %s", resp.Code, resp.Log)
	}
	return resp.Value, nil
}

func QuerySubspace(cdc *types.Codec, subspace []byte, storeName string) (res []types.KVPair, err error) {
	resRaw, err := Query(subspace, storeName, "subspace")
	if err != nil {
		return res, err
	}
	cdc.MustUnmarshalBinary(resRaw, &res)
	return
}
