package sync

import (
	"reflect"
	"sync"
	"time"

	"github.com/irisnet/iris-sync-server/model/store"
	"github.com/irisnet/iris-sync-server/model/store/document"
	"github.com/irisnet/iris-sync-server/module/logger"
	"github.com/irisnet/iris-sync-server/module/stake"
	"github.com/irisnet/iris-sync-server/util/constant"
	"github.com/irisnet/iris-sync-server/util/helper"
)

var (
	delay      = false
	methodName string
)

func handle(tx store.Docs, mutex sync.Mutex, funChains []func(tx store.Docs, mutex sync.Mutex)) {
	for _, fun := range funChains {
		fun(tx, mutex)
	}
}

// save Tx document into collection
func saveTx(tx store.Docs, mutex sync.Mutex) {
	methodName = "SaveTx: "
	txCollectionName := tx.Name()

	// save tx document into database
	storeTxDocFunc := func(doc store.Docs) {
		err := store.Save(doc)
		if err != nil {
			logger.Error.Printf("%v Save failed. doc is %+v, err is %v",
				methodName, doc, err.Error())
		}
	}

	switch txCollectionName {
	case document.CollectionNmCoinTx:
		storeTxDocFunc(tx)
		break
	case document.CollectionNmStakeTx:
		if !reflect.ValueOf(tx).FieldByName("Type").IsValid() {
			logger.Error.Printf("%v type which is field name of stake tx is missed\n", methodName)
			break
		}
		stakeType := constant.TxTypeStake + "/" + reflect.ValueOf(tx).FieldByName("Type").String()

		logger.Info.Printf("%v handle %v and type is %v",
			methodName, txCollectionName, stakeType)

		mutex.Lock()
		logger.Info.Printf("%v get lock\n", methodName)

		{
			switch stakeType {
			case stake.TypeTxDeclareCandidacy:
				stakeTxDeclareCandidacy, r := tx.(document.StakeTxDeclareCandidacy)
				if !r {
					logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
						methodName, stakeType)
					break
				}
				storeTxDocFunc(stakeTxDeclareCandidacy)

				cd, err := document.QueryCandidateByPubkey(stakeTxDeclareCandidacy.PubKey)

				candidate := document.Candidate{
					Address:     stakeTxDeclareCandidacy.From,
					PubKey:      stakeTxDeclareCandidacy.PubKey,
					Description: stakeTxDeclareCandidacy.Description,
				}
				// TODO: in further share not equal amount
				candidate.Shares += stakeTxDeclareCandidacy.Amount.Amount
				candidate.VotingPower += int64(stakeTxDeclareCandidacy.Amount.Amount)
				candidate.UpdateTime = stakeTxDeclareCandidacy.Time

				if err != nil && cd.Address == "" {
					logger.Warning.Printf("%v Replace candidate from %+v to %+v\n", methodName, cd, candidate)
				}

				store.SaveOrUpdate(candidate)
				break
			case stake.TypeTxDelegate:
				stakeTx, r := tx.(document.StakeTx)
				if !r {
					logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
						methodName, stakeType)
					break
				}
				storeTxDocFunc(stakeTx)

				candidate, err := document.QueryCandidateByPubkey(stakeTx.PubKey)
				// candidate is not exist
				if err != nil {
					logger.Warning.Printf("%v candidate is not exist while delegate, add = %s ,pub_key = %s\n",
						methodName, stakeTx.From, stakeTx.PubKey)
					candidate = document.Candidate{
						PubKey: stakeTx.PubKey,
					}
				}

				delegator, err := document.QueryDelegatorByAddressAndPubkey(stakeTx.From, stakeTx.PubKey)
				// delegator is not exist
				if err != nil {
					delegator = document.Delegator{
						Address: stakeTx.From,
						PubKey:  stakeTx.PubKey,
					}
				}
				// TODO: in further share not equal amount
				delegator.Shares += stakeTx.Amount.Amount
				delegator.UpdateTime = stakeTx.Time
				store.SaveOrUpdate(delegator)

				candidate.Shares += stakeTx.Amount.Amount
				candidate.VotingPower += int64(stakeTx.Amount.Amount)
				candidate.UpdateTime = stakeTx.Time
				store.SaveOrUpdate(candidate)
				break
			case stake.TypeTxUnbond:
				stakeTx, r := tx.(document.StakeTx)
				if !r {
					logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
						methodName, stakeType)
					break
				}
				storeTxDocFunc(stakeTx)

				delegator, err := document.QueryDelegatorByAddressAndPubkey(stakeTx.From, stakeTx.PubKey)
				// delegator is not exist
				if err != nil {
					logger.Warning.Printf("%v delegator is not exist while unBond,add = %s,pub_key=%s\n",
						methodName, stakeTx.From, stakeTx.PubKey)
					delegator = document.Delegator{
						Address: stakeTx.From,
						PubKey:  stakeTx.PubKey,
					}
				}
				delegator.Shares -= stakeTx.Amount.Amount
				delegator.UpdateTime = stakeTx.Time
				store.SaveOrUpdate(delegator)

				candidate, err2 := document.QueryCandidateByPubkey(stakeTx.PubKey)
				// candidate is not exist
				if err2 != nil {
					logger.Warning.Printf("%v candidate is not exist while unBond,add = %s,pub_key=%s\n",
						methodName, stakeTx.From, stakeTx.PubKey)
					candidate = document.Candidate{
						PubKey: stakeTx.PubKey,
					}
				}
				candidate.Shares -= stakeTx.Amount.Amount
				candidate.VotingPower -= int64(stakeTx.Amount.Amount)
				candidate.UpdateTime = stakeTx.Time
				store.SaveOrUpdate(candidate)
				break
			}
		}

		mutex.Unlock()
		logger.Info.Printf("%v method release lock\n", methodName)
	}
}

func saveOrUpdateAccount(tx store.Docs, mutex sync.Mutex) {
	var (
		address    string
		updateTime time.Time
		height     int64
		methodName = "SaveOrUpdateAccount: "
	)
	txCollectionName := tx.Name()
	logger.Info.Printf("Start %v\n", methodName)

	fun := func(address string, updateTime time.Time, height int64) {
		account := document.Account{
			Address: address,
			Time:    updateTime,
			Height:  height,
		}

		if err := store.Save(account); err != nil {
			logger.Trace.Printf("%v Record exists, account is %v, err is %s\n",
				methodName, account.Address, err.Error())
		}
	}

	// mutex.Lock()
	// logger.Info.Printf("%v get lock\n", methodName)

	{
		switch txCollectionName {
		case document.CollectionNmCoinTx:
			coinTx, r := tx.(document.CoinTx)
			if !r {
				logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
					methodName, txCollectionName)
				break
			}
			updateTime = coinTx.Time
			height = coinTx.Height

			fun(coinTx.From, updateTime, height)
			fun(coinTx.To, updateTime, height)
			break
		case document.CollectionNmStakeTx:
			if !reflect.ValueOf(tx).FieldByName("Type").IsValid() {
				logger.Error.Printf("%v type which is field name of stake tx is missed\n", methodName)
				break
			}
			stakeType := constant.TxTypeStake + "/" + reflect.ValueOf(tx).FieldByName("Type").String()

			logger.Info.Printf("%v handle %v and type is %v",
				methodName, txCollectionName, stakeType)

			switch stakeType {
			case stake.TypeTxDeclareCandidacy:
				stakeTxDeclareCandidacy, r := tx.(document.StakeTxDeclareCandidacy)
				if !r {
					logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
						methodName, stakeType)
					break
				}
				address = stakeTxDeclareCandidacy.From
				updateTime = stakeTxDeclareCandidacy.Time
				height = stakeTxDeclareCandidacy.Height
				break
			case stake.TypeTxEditCandidacy:
				break
			case stake.TypeTxDelegate, stake.TypeTxUnbond:
				stakeTx, r := tx.(document.StakeTx)
				if !r {
					logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
						methodName, stakeType)
					break
				}
				address = stakeTx.From
				updateTime = stakeTx.Time
				height = stakeTx.Height
				break
			}

			fun(address, updateTime, height)
		}
	}

	logger.Info.Printf("End %v\n", methodName)
	// mutex.Unlock()
	// logger.Info.Printf("%v release lock\n", methodName)
}

func updateAccountBalance(tx store.Docs, mutex sync.Mutex) {
	var (
		address    string
		methodName = "UpdateAccountBalance: "
	)
	txCollectionName := tx.Name()
	logger.Info.Printf("Start %v\n", methodName)

	fun := func(address string) {
		account, err := document.QueryAccount(address)
		if err != nil {
			logger.Error.Printf("%v updateAccountBalance failed, account is %v and err is %v",
				methodName, account, err.Error())
			return
		}

		// query balance of account
		ac := helper.QueryAccountBalance(address, delay)
		account.Amount = ac.Coins
		if err := store.Update(account); err != nil {
			logger.Error.Printf("%v account:[%q] balance update failed,%s\n",
				methodName, account.Address, err)
		}
	}

	// mutex.Lock()
	// logger.Info.Printf("%v get lock\n", methodName)
	{
		switch txCollectionName {
		case document.CollectionNmCoinTx:
			coinTx, r := tx.(document.CoinTx)
			if !r {
				logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
					methodName, txCollectionName)
				break
			}
			fun(coinTx.From)
			fun(coinTx.To)
			break
		case document.CollectionNmStakeTx:
			if !reflect.ValueOf(tx).FieldByName("Type").IsValid() {
				logger.Error.Printf("%v type which is field name of stake tx is missed\n", methodName)
				break
			}
			stakeType := constant.TxTypeStake + "/" + reflect.ValueOf(tx).FieldByName("Type").String()

			logger.Info.Printf("%v handle %v and type is %v",
				methodName, txCollectionName, stakeType)

			switch stakeType {
			case stake.TypeTxDeclareCandidacy:
				stakeTxDeclareCandidacy, r := tx.(document.StakeTxDeclareCandidacy)
				if !r {
					logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
						methodName, stakeType)
					break
				}
				address = stakeTxDeclareCandidacy.From
				break
			case stake.TypeTxEditCandidacy:
				break
			case stake.TypeTxDelegate, stake.TypeTxUnbond:
				stakeTx, r := tx.(document.StakeTx)
				if !r {
					logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
						methodName, stakeType)
					break
				}
				address = stakeTx.From
				break
			}
			fun(address)
			break
		}
	}
	
	logger.Info.Printf("End %v\n", methodName)

	// mutex.Unlock()
	// logger.Info.Println("updateAccountBalance method release lock")

}
