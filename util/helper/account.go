// This package is used for Query balance of account

package helper

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

// query account balance from sdk store
func QueryAccountInfo(address string) (store.Coins, uint64) {
	cdc := types.GetCodec()

	addr, err := types.AccAddressFromBech32(address)
	if err != nil {
		logger.Error("get addr from hex failed", logger.Any("err", err))
		return nil, 0
	}

	res, err := Query(types.AddressStoreKey(addr), "acc",
		constant.StoreDefaultEndPath)

	if err != nil {
		logger.Error("Query balance from tendermint failed", logger.Any("err", err))
		return nil, 0
	}

	// balance is empty
	if len(res) <= 0 {
		return nil, 0
	}

	decoder := types.GetAccountDecoder(cdc)
	account, err := decoder(res)
	if err != nil {
		logger.Error("decode account failed", logger.Any("err", err))
		return nil, 0
	}

	return types.ParseCoins(account.GetCoins().String()), account.GetAccountNumber()
}

func ValAddrToAccAddr(address string) (accAddr string) {
	valAddr, err := types.ValAddressFromBech32(address)
	if err != nil {
		logger.Error("ValAddressFromBech32 decode account failed", logger.String("address", address))
		return
	}

	return types.AccAddress(valAddr.Bytes()).String()
}
