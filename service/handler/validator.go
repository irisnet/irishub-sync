package handler

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/helper"
	"sort"
)

// compare validatorSet stored in irishub and validatorSet stored in db
// if dbCandidates not equal chainValidators, execute two step as follow.
// first, remove all validators stored in db
// second, store latest validators which query from sdk store into db
// note: this function isn't thread safe, should be invoked during watch block
//       not fast sync
func CompareAndUpdateValidators() {
	var (
		methodName = "CompareAndUpdateValidators"

		candidateModel document.Candidate
	)

	// get all validatorSets from db
	dbCandidates := candidateModel.QueryAll()

	// get all validatorSets from blockChain
	validators := helper.GetValidators()

	var chainValidators []document.Candidate
	for _, validator := range validators {
		// build validator document struct by stake.validator
		doc := BuildValidatorDocument(validator)
		chainValidators = append(chainValidators, doc)
	}

	// dbCandidates not equal chainValidators
	if compareValidators(dbCandidates, chainValidators) {
		// remove all data which stored in db
		if err := candidateModel.RemoveCandidates(); err != nil {
			logger.Error("RemoveCandidates err ", logger.String("method", methodName), logger.String("err", err.Error()))
		}

		updateValidatorsRank(&chainValidators)

		// store latest validators into db
		if err := candidateModel.SaveAll(chainValidators); err != nil {
			logger.Error("SaveAll", logger.String("method", methodName), logger.String("err", err.Error()))
		}
	}
}

func BuildValidatorDocument(v types.StakeValidator) document.Candidate {
	description := document.ValDescription{
		Moniker:  v.Description.Moniker,
		Identity: v.Description.Identity,
		Website:  v.Description.Website,
		Details:  v.Description.Details,
	}

	floatTokens := helper.ParseFloat(v.Tokens.String())
	floatDelegatorShares := helper.ParseFloat(v.DelegatorShares.String())
	pubKey, err := types.Bech32ifyValPub(v.ConsPubKey)
	if err != nil {
		logger.Error("Can't get validator pubKey", logger.String("pubKey", pubKey), logger.String("err", err.Error()))
	}
	doc := document.Candidate{
		Address:         v.OperatorAddr.String(),
		PubKey:          pubKey,
		PubKeyAddr:      v.ConsPubKey.Address().String(),
		Jailed:          v.Jailed,
		Tokens:          floatTokens,
		OriginalTokens:  helper.RoundString(v.Tokens.String(), 0),
		DelegatorShares: floatDelegatorShares,
		Description:     description,
		BondHeight:      v.BondHeight,
		Status:          types.BondStatusToString(v.Status),
	}

	doc.VotingPower = doc.Tokens

	return doc
}

func compareValidators(dbVals []document.Candidate, chainVals []document.Candidate) bool {
	//Candidate数量不一致
	if len(dbVals) != len(chainVals) {
		logger.Info("Candidate's member amount has changed")
		return true
	}

	chainValsMap := make(map[string]document.Candidate)
	for _, v := range chainVals {
		chainValsMap[v.PubKeyAddr] = v
	}

	for _, v := range dbVals {
		v1, ok := chainValsMap[v.PubKeyAddr]
		if !ok {
			logger.Info("Candidate's member has changed,removed",
				logger.String("dbValue", v.PubKeyAddr),
			)
			return true
		}

		if v.Tokens != v1.Tokens {
			logger.Info("Candidate's votingPower has changed",
				logger.String("validator", v.Address),
				logger.Float64("dbValue", v1.Tokens),
				logger.Float64("tmValue", v1.Tokens),
			)
			return true
		}

		if v.Jailed != v1.Jailed {
			logger.Info("Candidate's jailed status has changed",
				logger.String("validator", v.Address),
				logger.Bool("dbValue", v.Jailed),
				logger.Bool("tmValue", v1.Jailed),
			)
			return true
		}

		if v.Status != v1.Status {
			logger.Info("Candidate's status has changed",
				logger.String("validator", v.Address),
				logger.String("dbValue", v.Status),
				logger.String("tmValue", v1.Status),
			)
			return true
		}
	}
	logger.Info("Validators Set is not changed ")
	return false
}

func updateValidatorsRank(candidates *[]document.Candidate) {
	var historyModel document.ValidatorHistory
	vs := historyModel.QueryAll()
	if len(vs) == 0 {
		return
	}

	sort.Sort(CandidateWrapper{*candidates, func(p, q *document.Candidate) bool {
		return q.Tokens < p.Tokens // Tokens 递减排序
	}})

	var cMap = make(map[string]document.Candidate)

	for _, validator := range vs {
		cMap[validator.Address] = validator.Candidate
	}

	for index, candidate := range *candidates {
		lastCandidate := cMap[candidate.Address]
		var lift int
		if lastCandidate.Tokens > candidate.Tokens {
			lift = document.LiftDown
		} else if lastCandidate.Tokens < candidate.Tokens {
			lift = document.LiftUp
		} else {
			lift = document.LiftNotChange
		}
		rank := document.Rank{
			Number: index + 1,
			Lift:   lift,
		}
		(*candidates)[index].Rank = rank
	}
}

type CandidateWrapper struct {
	cs []document.Candidate
	by func(p, q *document.Candidate) bool
}

func (cw CandidateWrapper) Len() int { // 重写 Len() 方法
	return len(cw.cs)
}
func (cw CandidateWrapper) Swap(i, j int) { // 重写 Swap() 方法
	cw.cs[i], cw.cs[j] = cw.cs[j], cw.cs[i]
}
func (cw CandidateWrapper) Less(i, j int) bool { // 重写 Less() 方法
	return cw.by(&cw.cs[i], &cw.cs[j])
}
