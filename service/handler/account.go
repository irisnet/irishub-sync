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

func SaveOrUpdateAccountDelegationInfo(docTx document.CommonTx) {
	var (
		delegator    string
		accountModel document.Account
	)
	switch docTx.Type {
	case constant.TxTypeStakeDelegate, constant.TxTypeStakeBeginUnbonding, constant.TxTypeBeginRedelegate:
		delegator = docTx.From
	}
	if delegator == "" {
		return
	}
	delegations := helper.GetDelegations(delegator)
	delegation := store.Coin{
		Denom:  constant.IrisAttoUnit,
		Amount: helper.CalculateDelegatorDelegationTokens(delegations),
	}

	if err := accountModel.UpsertDelegationInfo(delegator, delegation, docTx.Height, docTx.Time.Unix()); err != nil {
		logger.Error("update account delegation info fail", logger.Int64("height", docTx.Height),
			logger.String("address", delegator), logger.String("err", err.Error()))
	}
}

func SaveOrUpdateAccountUnbondingDelegationInfo(docTx document.CommonTx) {
	var (
		delegator    string
		accountModel document.Account
	)
	switch docTx.Type {
	case constant.TxTypeStakeBeginUnbonding, constant.TxTypeBeginRedelegate:
		delegator = docTx.From
	}
	if delegator == "" {
		return
	}
	unbondingDelegations := helper.GetUnbondingDelegations(delegator)
	unbondingDelegation := store.Coin{
		Denom:  constant.IrisAttoUnit,
		Amount: helper.CalculateDelegatorUnbondingDelegationTokens(unbondingDelegations),
	}

	if err := accountModel.UpsertUnbondingDelegationInfo(delegator, unbondingDelegation, docTx.Height, docTx.Time.Unix()); err != nil {
		logger.Error("update account unbonding delegation info fail", logger.Int64("height", docTx.Height),
			logger.String("address", delegator), logger.String("err", err.Error()))
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
