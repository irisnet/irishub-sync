package sync

import (
	"github.com/irisnet/iris-sync-server/model/store"
	"github.com/irisnet/iris-sync-server/model/store/collection"
	"github.com/irisnet/iris-sync-server/module/logger"
	"github.com/irisnet/iris-sync-server/module/stake"
	"github.com/irisnet/iris-sync-server/util/helper"
)

var delay = false

func handle(tx store.Docs, funChains []func(tx store.Docs)) {
	for _, fun := range funChains {
		fun(tx)
	}
}

// save Tx document into collection
func saveTx(tx store.Docs) {
	store.Save(tx)

	if tx.Name() == collection.DocsNmStakeTx {
		stakeTx, _ := tx.(collection.StakeTx)

		switch stakeTx.Type {
		case stake.TypeTxUnbond:
			delegator, err := collection.QueryDelegatorByAddressAndPubkey(stakeTx.From, stakeTx.PubKey)
			if err != nil {
				logger.Info.Printf("error:delegator is lost,add = %s,pub_key=%s\n", stakeTx.From, stakeTx.PubKey)
				return
			}
			delegator.Shares -= stakeTx.Amount.Amount
			store.Update(delegator)

			candidate, err2 := collection.QueryCandidateByAddressAndPubkey(stakeTx.From, stakeTx.PubKey)
			if err2 != nil {
				logger.Info.Printf("error:candidate is lost,add = %s,pub_key=%s\n", stakeTx.From, stakeTx.PubKey)
				return
			}
			candidate.Shares -= stakeTx.Amount.Amount
			store.Update(candidate)

		case stake.TypeTxDelegate:
			de, err := collection.QueryDelegatorByAddressAndPubkey(stakeTx.From, stakeTx.PubKey)
			if err != nil {
				de = collection.Delegator{
					Address: stakeTx.From,
					PubKey:  stakeTx.PubKey,
				}
			}
			de.Shares += stakeTx.Amount.Amount
			store.SaveOrUpdate(de)
		case stake.TypeTxDeclareCandidacy:
			de, err := collection.QueryCandidateByAddressAndPubkey(stakeTx.From, stakeTx.PubKey)
			if err != nil {
				de = collection.Candidate{
					Address: stakeTx.From,
					PubKey:  stakeTx.PubKey,
				}
			}
			de.Shares += stakeTx.Amount.Amount
			store.SaveOrUpdate(de)
		}
	}
}

func saveOrUpdateAccount(tx store.Docs) {
	switch tx.Name() {
	case collection.DocsNmCoinTx:
		coinTx, _ := tx.(collection.CoinTx)
		fun := func(address string) {
			account := collection.Account{
				Address: address,
				Time:    coinTx.Time,
				Height:  coinTx.Height,
			}

			if err := store.SaveOrUpdate(account); err != nil {
				logger.Info.Printf("account:[%q] balance update failed,%s\n", account.Address, err)
			}
		}
		fun(coinTx.From)
		fun(coinTx.To)
	case collection.DocsNmStakeTx:
		stakeTx, _ := tx.(collection.StakeTx)
		fun := func(address string) {
			account := collection.Account{
				Address: address,
				Time:    stakeTx.Time,
				Height:  stakeTx.Height,
			}

			if err := store.SaveOrUpdate(account); err != nil {
				logger.Info.Printf("account:[%q] balance update failed,%s\n", account.Address, err)
			}
		}
		fun(stakeTx.From)
	}
}

func updateAccountBalance(tx store.Docs) {
	fun := func(address string) {
		account, _ := collection.QueryAccount(address)
		//查询账户余额
		ac := helper.QueryAccountBalance(address, delay)
		account.Amount = ac.Coins
		if err := store.Update(account); err != nil {
			logger.Info.Printf("account:[%q] balance update failed,%s\n", account.Address, err)
		}
	}
	switch tx.Name() {
	case collection.DocsNmCoinTx:
		coinTx, _ := tx.(collection.CoinTx)
		fun(coinTx.From)
		fun(coinTx.To)
	case collection.DocsNmStakeTx:
		stakeTx, _ := tx.(collection.StakeTx)
		fun(stakeTx.From)
	case collection.DocsNmAccount:
		account, _ := tx.(collection.Account)
		ac := helper.QueryAccountBalance(account.Address, delay)
		account.Amount = ac.Coins
		if err := store.Update(account); err != nil {
			logger.Info.Printf("account:[%q] balance update failed,%s\n", account.Address, err)
		}
	}

}
