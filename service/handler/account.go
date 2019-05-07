package handler

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
)

func SaveOrUpdateAccountBalanceInfo(accounts []string, height, timestamp int64) {
	var (
		accountModel document.Account
	)
	if len(accounts) == 0 {
		return
	}

	for _, v := range accounts {
		coins, accountNumber := helper.QueryAccountInfo(v)
		coinIris := getCoinIrisFromCoins(coins)

		if err := accountModel.UpsertBalanceInfo(v, coinIris, accountNumber, height, timestamp); err != nil {
			logger.Error("update account balance info fail", logger.Int64("height", height),
				logger.String("address", v), logger.String("err", err.Error()))
		}
	}
}

func SaveOrUpdateAccountDelegationInfo() {

}

func SaveOrUpdateAccountUnbondingDelegationInfo() {

}

func getCoinIrisFromCoins(coins store.Coins) store.Coin {
	if len(coins) > 0 {
		for _, v := range coins {
			if v.Denom == constant.IrisAttoUnit {
				return store.Coin{
					Denom:  v.Denom,
					Amount: v.Amount,
				}
			}
		}
	}
	return store.Coin{}
}
