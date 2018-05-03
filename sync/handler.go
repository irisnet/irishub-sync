package sync

import (
	"github.com/irisnet/iris-sync-server/model/store"
	"github.com/irisnet/iris-sync-server/module/logger"
	"github.com/irisnet/iris-sync-server/module/stake"
	"github.com/irisnet/iris-sync-server/util/helper"
	"github.com/irisnet/iris-sync-server/model/store/document"
	"reflect"
	"github.com/irisnet/iris-sync-server/util/constant"
)

var delay = false

func handle(tx store.Docs, funChains []func(tx store.Docs)) {
	for _, fun := range funChains {
		fun(tx)
	}
}

// save Tx document into collection
func saveTx(tx store.Docs) {
	err := store.Save(tx)

	logger.Info.Printf("collection name is %s\n", tx.Name())
	if err != nil {
		logger.Info.Println(err)
	}

	// TODO: Thread safe?
	if tx.Name() == document.CollectionNmStakeTx {
		if !reflect.ValueOf(tx).FieldByName("Type").IsValid() {
			logger.Error.Println("type which is field name of stake tx is missed")
			return
		}
		stakeType := constant.TxTypeStake + "/" + reflect.ValueOf(tx).FieldByName("Type").String()

		switch stakeType {
		case stake.TypeTxDeclareCandidacy:
			stakeTxDeclareCandidacy, _ := tx.(document.StakeTxDeclareCandidacy)
			candidate, err := document.QueryCandidateByPubkey(stakeTxDeclareCandidacy.PubKey)
			// new candidate
			if err != nil {
				candidate = document.Candidate {
					Address:     stakeTxDeclareCandidacy.From,
					PubKey:      stakeTxDeclareCandidacy.PubKey,
					Description: stakeTxDeclareCandidacy.Description,
				}
			}
			// TODO: in further share not equal amount
			candidate.Shares += stakeTxDeclareCandidacy.Amount.Amount
			candidate.VotingPower += uint64(stakeTxDeclareCandidacy.Amount.Amount)
			store.SaveOrUpdate(candidate)
			break
		case stake.TypeTxDelegate:
			stakeTx, _ := tx.(document.StakeTx)
			candidate, err := document.QueryCandidateByPubkey(stakeTx.PubKey)
			// candidate is not exist
			if err != nil {
				logger.Error.Printf("candidate is not exist while delegate, add = %s ,pub_key = %s\n",
					stakeTx.From, stakeTx.PubKey)
				return
			}

			delegator, err := document.QueryDelegatorByAddressAndPubkey(stakeTx.From, stakeTx.PubKey)
			// delegator is not exist
			if err != nil {
				delegator = document.Delegator{
					Address: stakeTx.From,
					PubKey: stakeTx.PubKey,
				}
			}
			// TODO: in further share not equal amount
			delegator.Shares += stakeTx.Amount.Amount
			store.SaveOrUpdate(delegator)

			candidate.Shares += stakeTx.Amount.Amount
			candidate.VotingPower += uint64(stakeTx.Amount.Amount)
			store.SaveOrUpdate(candidate)
			break
		case stake.TypeTxUnbond:
			stakeTx, _ := tx.(document.StakeTx)
			delegator, err := document.QueryDelegatorByAddressAndPubkey(stakeTx.From, stakeTx.PubKey)
			// delegator is not exist
			if err != nil {
				logger.Info.Printf("delegator is not exist while unBond,add = %s,pub_key=%s\n",
					stakeTx.From, stakeTx.PubKey)
				return
			}
			delegator.Shares -= stakeTx.Amount.Amount
			store.Update(delegator)

			candidate, err2 := document.QueryCandidateByPubkey(stakeTx.PubKey)
			// candidate is not exist
			if err2 != nil {
				logger.Info.Printf("candidate is not exist while unBond,add = %s,pub_key=%s\n",
					stakeTx.From, stakeTx.PubKey)
				return
			}
			candidate.Shares -= stakeTx.Amount.Amount
			candidate.VotingPower -= uint64(stakeTx.Amount.Amount)
			store.Update(candidate)
			break
		}

	}
}

func saveOrUpdateAccount(tx store.Docs) {
	switch tx.Name() {
	case document.CollectionNmCoinTx:
		coinTx, _ := tx.(document.CoinTx)
		fun := func(address string) {
			account := document.Account{
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
	case document.CollectionNmStakeTx:
		stakeTx, _ := tx.(document.StakeTx)
		fun := func(address string) {
			account := document.Account{
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
		account, _ := document.QueryAccount(address)
		//查询账户余额
		ac := helper.QueryAccountBalance(address, delay)
		account.Amount = ac.Coins
		if err := store.Update(account); err != nil {
			logger.Info.Printf("account:[%q] balance update failed,%s\n", account.Address, err)
		}
	}
	switch tx.Name() {
	case document.CollectionNmCoinTx:
		coinTx, _ := tx.(document.CoinTx)
		fun(coinTx.From)
		fun(coinTx.To)
	case document.CollectionNmStakeTx:
		stakeTx, _ := tx.(document.StakeTx)
		fun(stakeTx.From)
	case document.CollectionNmAccount:
		account, _ := tx.(document.Account)
		ac := helper.QueryAccountBalance(account.Address, delay)
		account.Amount = ac.Coins
		if err := store.Update(account); err != nil {
			logger.Info.Printf("account:[%q] balance update failed,%s\n", account.Address, err)
		}
	}

}
