// This package is used for Query balance of account

package helper

import (
	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/store"

	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/pkg/errors"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// query account balance from sdk store
func QueryAccountBalance(address string) store.Coins {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		logger.Error.Printf("get addr from hex failed, %+v\n", err)
		return nil
	}

	res, err := Query(auth.AddressStoreKey(addr), "acc",
		constant.StoreDefaultEndPath)

	if err != nil {
		logger.Error.Printf("Query balance from tendermint failed, %+v\n", err)
		return nil
	}

	// balance is empty
	if len(res) <= 0 {
		return nil
	}

	decoder := authcmd.GetAccountDecoder(codec.Cdc)
	account, err := decoder(res)
	if err != nil {
		logger.Error.Printf("decode account failed, %+v\n", err)
		return nil
	}

	return BuildCoins(account.GetCoins())
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
