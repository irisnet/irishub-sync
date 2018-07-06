package handler

import (
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"sync"
)

// save Tx document into collection
func SaveTx(docTx store.Docs, mutex sync.Mutex) {
	var (
		methodName = "SaveTx: "
	)

	// save docTx document into database
	storeTxDocFunc := func(doc store.Docs) {
		err := store.Save(doc)
		if err != nil {
			logger.Error.Printf("%v Save failed. doc is %+v, err is %v",
				methodName, doc, err.Error())
		}
	}

	// save common docTx document
	saveCommonTx := func(commonTx document.CommonTx) {
		err := store.Save(commonTx)
		if err != nil {
			logger.Error.Printf("%v Save commonTx failed. doc is %+v, err is %v",
				methodName, commonTx, err.Error())
		}
	}

	txType := GetTxType(docTx)
	if txType == "" {
		logger.Error.Printf("%v get docTx type failed, docTx is %v\n",
			methodName, docTx)
		return
	}

	saveCommonTx(buildCommonTxData(docTx, txType))

	switch txType {
	case constant.TxTypeBank:
		break
	case constant.TxTypeStakeCreate:
		docTx, r := docTx.(document.StakeTxDeclareCandidacy)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		storeTxDocFunc(docTx)

		mutex.Lock()
		logger.Info.Printf("%v saveOrUpdate cndidates get lock\n", methodName)

		cd, err := document.QueryCandidateByAddress(docTx.ValidatorAddr)

		candidate := document.Candidate{
			Address: docTx.ValidatorAddr,
			PubKey:  docTx.PubKey,
		}

		if err == nil && cd.PubKey == "" {
			// candidate exist
			logger.Warning.Printf("%v Replace candidate from %+v to %+v\n", methodName, cd, candidate)
			// TODO: in further share not equal amount
			candidate.Shares = cd.Shares + docTx.Amount.Amount
			candidate.VotingPower = int64(candidate.Shares)

			// description of candidate is empty
			if cd.Description.Moniker == "" {
				candidate.Description = docTx.Description
			}
		} else {
			// candidate not exist
			candidate.Shares = docTx.Amount.Amount
			candidate.VotingPower = int64(candidate.Shares)
			candidate.Description = docTx.Description
		}
		candidate.UpdateTime = docTx.Time

		store.SaveOrUpdate(candidate)

		mutex.Unlock()
		logger.Info.Printf("%v saveOrUpdate cndidates release lock\n", methodName)
		break

	case constant.TxTypeStakeEdit:
		docTx, r := docTx.(document.StakeTxEditCandidacy)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		storeTxDocFunc(docTx)

		mutex.Lock()
		logger.Info.Printf("%v saveOrUpdate cndidates get lock\n", methodName)

		cd, err := document.QueryCandidateByAddress(docTx.ValidatorAddr)

		var candidate document.Candidate

		if err != nil {
			// candidate not exist
			candidate = document.Candidate{
				Address:     docTx.ValidatorAddr,
				Description: docTx.Description,
			}
		} else {
			// candidate exist
			cd.Description = docTx.Description
			candidate = cd
		}

		store.SaveOrUpdate(candidate)

		mutex.Unlock()
		logger.Info.Printf("%v saveOrUpdate cndidates release lock\n", methodName)
		break
	case constant.TxTypeStakeDelegate:
		docTx, r := docTx.(document.StakeTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		storeTxDocFunc(docTx)

		mutex.Lock()
		logger.Info.Printf("%v saveOrUpdate docTx type %v get lock\n",
			methodName, txType)

		candidate, err := document.QueryCandidateByAddress(docTx.ValidatorAddr)
		// candidate is not exist
		if err != nil {
			logger.Warning.Printf("%v candidate is not exist while delegate, addrdelegator = %s ,addrvalidator = %s\n",
				methodName, docTx.DelegatorAddr, docTx.ValidatorAddr)
			candidate = document.Candidate{
				Address: docTx.ValidatorAddr,
			}
		}

		delegator, err := document.QueryDelegatorByAddressAndValAddr(docTx.DelegatorAddr, docTx.ValidatorAddr)
		// delegator is not exist
		if err != nil {
			delegator = document.Delegator{
				Address:       docTx.DelegatorAddr,
				ValidatorAddr: docTx.ValidatorAddr,
			}
		}
		// TODO: in further share not equal amount
		delegator.Shares += docTx.Amount.Amount
		delegator.UpdateTime = docTx.Time
		store.SaveOrUpdate(delegator)

		candidate.Shares += docTx.Amount.Amount
		candidate.VotingPower += int64(docTx.Amount.Amount)
		candidate.UpdateTime = docTx.Time
		store.SaveOrUpdate(candidate)

		mutex.Unlock()
		logger.Info.Printf("%v saveOrUpdate docTx type %v release lock\n",
			methodName, txType)
		break

	case constant.TxTypeStakeUnbond:
		docTx, r := docTx.(document.StakeTx)
		if !r {
			logger.Error.Printf("%v get docuemnt from docTx failed. docTx type is %v\n",
				methodName, txType)
			break
		}
		storeTxDocFunc(docTx)

		mutex.Lock()
		logger.Info.Printf("%v saveOrUpdate docTx type %v get lock\n",
			methodName, txType)

		delegator, err := document.QueryDelegatorByAddressAndValAddr(docTx.DelegatorAddr, docTx.ValidatorAddr)
		// delegator is not exist
		if err != nil {
			logger.Warning.Printf("%v delegator is not exist while unBond,add = %s,valAddr=%s\n",
				methodName, docTx.DelegatorAddr, docTx.ValidatorAddr)
			delegator = document.Delegator{
				Address:       docTx.DelegatorAddr,
				ValidatorAddr: docTx.ValidatorAddr,
			}
		}
		delegator.Shares -= docTx.Amount.Amount
		delegator.UpdateTime = docTx.Time
		store.SaveOrUpdate(delegator)

		candidate, err2 := document.QueryCandidateByAddress(docTx.ValidatorAddr)
		// candidate is not exist
		if err2 != nil {
			logger.Warning.Printf("%v candidate is not exist while unBond,add = %s,valAddr=%s\n",
				methodName, docTx.DelegatorAddr, docTx.ValidatorAddr)
			candidate = document.Candidate{
				Address: docTx.ValidatorAddr,
			}
		}
		candidate.Shares -= docTx.Amount.Amount
		candidate.VotingPower -= int64(docTx.Amount.Amount)
		candidate.UpdateTime = docTx.Time
		store.SaveOrUpdate(candidate)

		mutex.Unlock()
		logger.Info.Printf("%v saveOrUpdate docTx type %v release lock\n",
			methodName, txType)
		break
	}
}

// build common tx data through parse tx
func buildCommonTxData(docTx store.Docs, txType string) document.CommonTx {
	var commonTx document.CommonTx

	if txType == "" {
		txType = GetTxType(docTx)
	}
	switch txType {
	case constant.TxTypeBank:
		doc := docTx.(document.CommonTx)
		commonTx = document.CommonTx{
			TxHash: doc.TxHash,
			Time:   doc.Time,
			Height: doc.Height,
			From:   doc.From,
			To:     doc.To,
			Amount: doc.Amount,
			Type:   doc.Type,
			Fee:    doc.Fee,
			Status: doc.Status,
		}
		break
	case constant.TxTypeStakeCreate:
		doc := docTx.(document.StakeTxDeclareCandidacy)
		commonTx = document.CommonTx{
			TxHash: doc.TxHash,
			Time:   doc.Time,
			Height: doc.Height,
			From:   doc.ValidatorAddr,
			To:     "",
			Amount: []store.Coin{doc.Amount},
			Type:   doc.Type,
			Fee:    doc.Fee,
			Status: doc.Status,
		}
		break
	case constant.TxTypeStakeEdit:
		doc := docTx.(document.StakeTxEditCandidacy)
		commonTx = document.CommonTx{
			TxHash: doc.TxHash,
			Time:   doc.Time,
			Height: doc.Height,
			From:   doc.ValidatorAddr,
			To:     "",
			Amount: []store.Coin{doc.Amount},
			Type:   doc.Type,
			Fee:    doc.Fee,
			Status: doc.Status,
		}
		break
	case constant.TxTypeStakeDelegate, constant.TxTypeStakeUnbond:
		doc := docTx.(document.StakeTx)
		commonTx = document.CommonTx{
			TxHash: doc.TxHash,
			Time:   doc.Time,
			Height: doc.Height,
			From:   doc.ValidatorAddr,
			To:     doc.DelegatorAddr,
			Amount: []store.Coin{doc.Amount},
			Type:   doc.Type,
			Fee:    doc.Fee,
			Status: doc.Status,
		}
		break
	}

	return commonTx
}
