// This package is used for query balance of account

package helper

import (
	"github.com/irisnet/irishub-sync/model/store"
	"github.com/irisnet/irishub-sync/module/codec"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	cmn "github.com/tendermint/tmlibs/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"fmt"
	"github.com/irisnet/irishub-sync/module/logger"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
)

func QueryAccountBalance(address string) store.Coins {
	addr, err := sdk.GetValAddressHex(address)
	if err != nil {
		logger.Error.Printf("get addr from hex failed, %+v\n", err)
		return nil
	}

	res, err := query(auth.AddressStoreKey(addr), "acc", "key")

	if err != nil {
		logger.Error.Printf("query balance from tendermint failed, %+v\n", err)
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
func query(key cmn.HexBytes, storeName string, endPath string) (res []byte, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	rpcClient := GetClient().Client
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

