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
	"github.com/cosmos/cosmos-sdk/modules/coin"
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

	// save tx document into database
	storeTxDocFunc := func(doc store.Docs) {
		err := store.Save(doc)
		if err != nil {
			logger.Error.Printf("%v Save failed. doc is %+v, err is %v",
				methodName, doc, err.Error())
		}
	}
	
	// save common tx document
	saveCommonTx := func(commonTx document.CommonTx) {
		err := store.Save(commonTx)
		if err != nil {
			logger.Error.Printf("%v Save commonTx failed. doc is %+v, err is %v",
				methodName, commonTx, err.Error())
		}
	}
	
	txType := GetTxType(tx)
	if txType == "" {
		logger.Error.Printf("%v get tx type failed, tx is %v\n",
			methodName, tx)
		return
	}
	
	saveCommonTx(buildCommonTxData(tx, txType))
	
	switch txType {
	case constant.TxTypeCoin:
		storeTxDocFunc(tx)
		break
	case stake.TypeTxDeclareCandidacy:
		stakeTxDeclareCandidacy, r := tx.(document.StakeTxDeclareCandidacy)
		if !r {
			logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
				methodName, txType)
			break
		}
		storeTxDocFunc(stakeTxDeclareCandidacy)
		
		mutex.Lock()
		logger.Info.Printf("%v saveOrUpdate cndidates get lock\n", methodName)
		
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
		
		mutex.Unlock()
		logger.Info.Printf("%v saveOrUpdate cndidates release lock\n", methodName)
		break
	
	case stake.TypeTxEditCandidacy:
		break
	case stake.TypeTxDelegate:
		stakeTx, r := tx.(document.StakeTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
				methodName, txType)
			break
		}
		storeTxDocFunc(stakeTx)
		
		mutex.Lock()
		logger.Info.Printf("%v saveOrUpdate tx type %v get lock\n",
			methodName, txType)
		
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
		
		mutex.Unlock()
		logger.Info.Printf("%v saveOrUpdate tx type %v release lock\n",
			methodName, txType)
		break
	
	case stake.TypeTxUnbond:
		stakeTx, r := tx.(document.StakeTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
				methodName, txType)
			break
		}
		storeTxDocFunc(stakeTx)
		
		mutex.Lock()
		logger.Info.Printf("%v saveOrUpdate tx type %v get lock\n",
			methodName, txType)
		
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
		
		mutex.Unlock()
		logger.Info.Printf("%v saveOrUpdate tx type %v release lock\n",
			methodName, txType)
		break
	}
}

// save account
func saveOrUpdateAccount(tx store.Docs, mutex sync.Mutex) {
	var (
		address    string
		updateTime time.Time
		height     int64
		methodName = "SaveOrUpdateAccount: "
	)
	logger.Info.Printf("Start %v\n", methodName)

	// save account
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
	
	txType := GetTxType(tx)
	if txType == "" {
		logger.Error.Printf("%v get tx type failed, tx is %v\n",
			methodName, tx)
		return
	}
	
	switch txType {
	case constant.TxTypeCoin:
		coinTx, r := tx.(document.CoinTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
				methodName, txType)
			break
		}
		updateTime = coinTx.Time
		height = coinTx.Height
		
		fun(coinTx.From, updateTime, height)
		fun(coinTx.To, updateTime, height)
		break
	case stake.TypeTxDeclareCandidacy:
		stakeTxDeclareCandidacy, r := tx.(document.StakeTxDeclareCandidacy)
		if !r {
			logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
				methodName, txType)
			break
		}
		address = stakeTxDeclareCandidacy.From
		updateTime = stakeTxDeclareCandidacy.Time
		height = stakeTxDeclareCandidacy.Height
		
		fun(address, updateTime, height)
		break
	case stake.TypeTxEditCandidacy:
		break
	case stake.TypeTxDelegate, stake.TypeTxUnbond:
		stakeTx, r := tx.(document.StakeTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
				methodName, txType)
			break
		}
		address = stakeTx.From
		updateTime = stakeTx.Time
		height = stakeTx.Height
		
		fun(address, updateTime, height)
		break
	}
	
	logger.Info.Printf("End %v\n", methodName)
}

// update account balance
func updateAccountBalance(tx store.Docs, mutex sync.Mutex) {
	var (
		address    string
		methodName = "UpdateAccountBalance: "
	)
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
	
	txType := GetTxType(tx)
	if txType == "" {
		logger.Error.Printf("%v get tx type failed, tx is %v\n",
			methodName, tx)
		return
	}
	
	switch txType {
	case constant.TxTypeCoin:
		coinTx, r := tx.(document.CoinTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
				methodName, txType)
			break
		}
		fun(coinTx.From)
		fun(coinTx.To)
		break
	case stake.TypeTxDeclareCandidacy:
		stakeTxDeclareCandidacy, r := tx.(document.StakeTxDeclareCandidacy)
		if !r {
			logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
				methodName, txType)
			break
		}
		address = stakeTxDeclareCandidacy.From
		fun(address)
		break
	case stake.TypeTxEditCandidacy:
		break
	case stake.TypeTxDelegate, stake.TypeTxUnbond:
		stakeTx, r := tx.(document.StakeTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from tx failed. tx type is %v\n",
				methodName, txType)
			break
		}
		address = stakeTx.From
		fun(address)
		break
	}
	
	logger.Info.Printf("End %v\n", methodName)
}

// build common tx data through parse tx
func buildCommonTxData(tx store.Docs, txType string) document.CommonTx {
	var commonTx document.CommonTx
	
	if txType == "" {
		txType = GetTxType(tx)
	}
	switch txType {
	// tx coin
	case constant.TxTypeCoin:
		txCoin := tx.(document.CoinTx)
		commonTx = document.CommonTx{
			TxHash: txCoin.TxHash,
			Time: txCoin.Time,
			Height: txCoin.Height,
			From: txCoin.From,
			To: txCoin.To,
			Amount: txCoin.Amount,
			Type: txType,
		}
		break
	case stake.TypeTxDeclareCandidacy:
		txDeclare := tx.(document.StakeTxDeclareCandidacy)
		commonTx = document.CommonTx{
			TxHash: txDeclare.TxHash,
			Time: txDeclare.Time,
			Height: txDeclare.Height,
			From: txDeclare.From,
			To: txDeclare.PubKey,
			Amount: []coin.Coin{txDeclare.Amount},
			Type: txDeclare.Type,
		}
		break
	case stake.TypeTxEditCandidacy:
		break
	case stake.TypeTxDelegate, stake.TypeTxUnbond:
		txStake := tx.(document.StakeTx)
		commonTx = document.CommonTx{
			TxHash: txStake.TxHash,
			Time: txStake.Time,
			Height: txStake.Height,
			From: txStake.From,
			To: txStake.PubKey,
			Amount: []coin.Coin{txStake.Amount},
			Type: txStake.Type,
		}
		break
	}
	
	return commonTx
}

// get tx type
func GetTxType(tx store.Docs) string {
	txCollectionName := tx.Name()
	var txType string
	
	switch txCollectionName {
	case document.CollectionNmCoinTx:
		txType = constant.TxTypeCoin
		break
	case document.CollectionNmStakeTx:
		if !reflect.ValueOf(tx).FieldByName("Type").IsValid() {
			logger.Error.Printf("%v type which is field name of stake tx is missed\n", methodName)
			break
		}
		stakeType := constant.TxTypeStake + "/" + reflect.ValueOf(tx).FieldByName("Type").String()
		
		switch stakeType {
		case stake.TypeTxDeclareCandidacy:
			txType = stake.TypeTxDeclareCandidacy
			break
		case stake.TypeTxEditCandidacy:
			txType = stake.TypeTxEditCandidacy
			break
		case stake.TypeTxDelegate:
			txType = stake.TypeTxDelegate
			break
		case stake.TypeTxUnbond:
			txType = stake.TypeTxUnbond
			break
		}
		break
	}
	
	return txType
}