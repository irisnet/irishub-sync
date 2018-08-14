package helper

import (
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cosmos/cosmos-sdk/x/stake/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/module/codec"

	"fmt"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/pkg/errors"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// get validator
func GetValidator(valAddr string) (stake.Validator, error) {
	var (
		validatorAddr sdk.AccAddress
		err           error
		res           stake.Validator
	)

	validatorAddr, err = sdk.AccAddressFromBech32(valAddr)

	resRaw, err := Query(stake.GetValidatorKey(validatorAddr), constant.StoreNameStake, constant.StoreDefaultEndPath)
	if err != nil || resRaw == nil {
		return res, err
	}

	res = types.MustUnmarshalValidator(codec.Cdc, validatorAddr, resRaw)

	return res, err
}

// get delegation
func GetDelegation(delAddr, valAddr string) (stake.Delegation, error) {
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
func GetUnbondingDelegation(delAddr, valAddr string) (stake.UnbondingDelegation, error) {
	var (
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.AccAddress
		err           error

		res stake.UnbondingDelegation
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
	key := stake.GetUBDKey(delegatorAddr, validatorAddr)

	resRaw, err := Query(key, constant.StoreNameStake, constant.StoreDefaultEndPath)

	if err != nil || resRaw == nil {
		return res, err
	}

	res = types.MustUnmarshalUBD(cdc, key, resRaw)

	return res, nil
}

// Query from Tendermint with the provided storename and path
func Query(key cmn.HexBytes, storeName string, endPath string) (res []byte, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	client := GetClient()
	defer client.Release()

	rpcClient := client.Client
	if err != nil {
		return res, err
	}

	opts := rpcclient.ABCIQueryOptions{
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
