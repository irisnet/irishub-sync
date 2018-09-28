// This package is used for Query balance of account

package helper

import (
	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

// query account balance from sdk store
func QueryAccountBalance(address string) store.Coins {
	addr, err := types.AccAddressFromBech32(address)
	if err != nil {
		logger.Error("get addr from hex failed", logger.String("err", err.Error()))
		return nil
	}

	res, err := Query(types.AddressStoreKey(addr), "acc",
		constant.StoreDefaultEndPath)

	if err != nil {
		logger.Error("Query balance from tendermint failed", logger.String("err", err.Error()))
		return nil
	}

	// balance is empty
	if len(res) <= 0 {
		return nil
	}

	decoder := types.GetAccountDecoder(codec.Cdc)
	account, err := decoder(res)
	if err != nil {
		logger.Error("decode account failed", logger.String("err", err.Error()))
		return nil
	}

	return types.BuildCoins(account.GetCoins())
}
