package handler

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
)

// update account info
func UpdateAccountInfo(accounts []string, blockTime int64) {
	if len(accounts) == 0 {
		return
	}
	for _, v := range accounts {
		coins, accountNumber := helper.QueryAccountInfo(v)
		coinIris := getCoinIrisFromCoins(coins)
		account := document.Account{
			Address:          v,
			AccountNumber:    accountNumber,
			CoinIris:         coinIris,
			CoinIrisUpdateAt: blockTime,
		}

		if err := store.Upsert(account); err != nil {
			logger.Error("upsert account info fail", logger.String("addr", v),
				logger.String("err", err.Error()))
		}
	}
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
