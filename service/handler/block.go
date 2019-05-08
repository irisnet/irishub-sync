package handler

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/helper"
	"strings"
)

var (
	assetDetailTriggers = map[string]bool{
		"stakeEndBlocker":   true,
		"slashBeginBlocker": true,
		"slashEndBlocker":   true,
		"govEndBlocker":     true,
	}

	bech32AccountAddrPrefix = itypes.Bech32AccountAddrPrefix
)

const (
	triggerTxHashLength = 64
	separator           = "::" // tag value separator
	unDelegationSubject = "Undelegation"
)

func ParseBlock(meta *types.BlockMeta, block *types.Block, validators []*types.Validator) document.Block {
	cdc := types.GetCodec()

	hexFunc := func(bytes []byte) string {
		return helper.BuildHex(bytes)
	}

	docBlock := document.Block{
		Height: meta.Header.Height,
		Hash:   hexFunc(meta.BlockID.Hash),
		Time:   meta.Header.Time,
		NumTxs: meta.Header.NumTxs,
	}

	lastBlockId := document.BlockID{
		Hash: hexFunc(meta.Header.LastBlockID.Hash),
		PartsHeader: document.PartSetHeader{
			Total: meta.Header.LastBlockID.PartsHeader.Total,
			Hash:  hexFunc(meta.Header.LastBlockID.PartsHeader.Hash),
		},
	}

	// blockMeta
	blockMeta := document.BlockMeta{
		BlockID: document.BlockID{
			Hash: hexFunc(meta.BlockID.Hash),
			PartsHeader: document.PartSetHeader{
				Total: meta.BlockID.PartsHeader.Total,
				Hash:  hexFunc(meta.BlockID.PartsHeader.Hash),
			},
		},
		Header: document.Header{
			ChainID:         meta.Header.ChainID,
			Height:          meta.Header.Height,
			Time:            meta.Header.Time,
			NumTxs:          meta.Header.NumTxs,
			LastBlockID:     lastBlockId,
			TotalTxs:        meta.Header.TotalTxs,
			LastCommitHash:  hexFunc(meta.Header.LastCommitHash),
			DataHash:        hexFunc(meta.Header.DataHash),
			ValidatorsHash:  hexFunc(meta.Header.ValidatorsHash),
			ConsensusHash:   hexFunc(meta.Header.ConsensusHash),
			AppHash:         hexFunc(meta.Header.AppHash),
			LastResultsHash: hexFunc(meta.Header.LastResultsHash),
			EvidenceHash:    hexFunc(meta.Header.EvidenceHash),
		},
	}

	// block
	var (
		preCommits []document.Vote
	)

	if len(block.LastCommit.Precommits) > 0 {
		for _, v := range block.LastCommit.Precommits {
			if v != nil {
				var sig document.Signature
				out, _ := cdc.MarshalJSON(v.Signature)
				json.Unmarshal(out, &sig)
				preCommit := document.Vote{
					ValidatorAddress: v.ValidatorAddress.String(),
					ValidatorIndex:   v.ValidatorIndex,
					Height:           v.Height,
					Round:            v.Round,
					Timestamp:        v.Timestamp,
					Type:             byte(v.Type),
					BlockID:          lastBlockId,
					Signature:        sig,
				}
				preCommits = append(preCommits, preCommit)
			}
		}
	}

	blockContent := document.BlockContent{
		LastCommit: document.Commit{
			BlockID:    lastBlockId,
			Precommits: preCommits,
		},
	}

	// validators
	var vals []document.Validator
	if len(validators) > 0 {
		for _, v := range validators {
			validator := document.Validator{
				Address:     v.Address.String(),
				VotingPower: v.VotingPower,
				PubKey:      hexFunc(v.PubKey.Bytes()),
			}
			vals = append(vals, validator)
		}
	}

	docBlock.Meta = blockMeta
	docBlock.Block = blockContent
	docBlock.Validators = vals
	docBlock.Result = parseBlockResult(docBlock.Height)

	// save or update account balance info and unbonding delegation info by parse block coin flow
	accsBalanceNeedUpdated, accsUnbondingDelegationNeedUpdated := getAccountsFromCoinFlow(
		docBlock.Result.EndBlock.Tags, docBlock.Height)
	SaveOrUpdateAccountBalanceInfo(accsBalanceNeedUpdated, docBlock.Height, docBlock.Time.Unix())
	SaveOrUpdateAccountUnbondingDelegationInfo(accsUnbondingDelegationNeedUpdated, docBlock.Height, docBlock.Time.Unix())

	return docBlock
}

func parseBlockResult(height int64) (res document.BlockResults) {
	client := helper.GetClient()
	defer client.Release()

	result, err := client.BlockResults(&height)
	if err != nil {
		logger.Error("EndBlocker error", logger.Any("err", err))
	}

	var deliverTxRes []document.ResponseDeliverTx
	for _, tx := range result.Results.DeliverTx {
		deliverTxRes = append(deliverTxRes, document.ResponseDeliverTx{
			Code:      tx.Code,
			Data:      string(tx.Data),
			Log:       tx.Log,
			GasWanted: tx.GasWanted,
			GasUsed:   tx.GasUsed,
			Tags:      parseTags(tx.Tags),
		})
	}

	res.DeliverTx = deliverTxRes

	var validatorUpdates []document.ValidatorUpdate
	for _, tx := range result.Results.EndBlock.ValidatorUpdates {
		validatorUpdates = append(validatorUpdates, document.ValidatorUpdate{
			PubKey: tx.PubKey.String(),
			Power:  tx.Power,
		})
	}

	var consensusParamUpdates document.ConsensusParams
	var tmConsensusParamUpdates = result.Results.EndBlock.ConsensusParamUpdates
	if tmConsensusParamUpdates != nil {
		if tmConsensusParamUpdates.Validator != nil {
			consensusParamUpdates.Validator = document.ValidatorParams{
				PubKeyTypes: tmConsensusParamUpdates.Validator.PubKeyTypes,
			}
		}
		if tmConsensusParamUpdates.BlockSize != nil {
			consensusParamUpdates.BlockSize = document.BlockSizeParams{
				MaxBytes: tmConsensusParamUpdates.BlockSize.MaxBytes,
				MaxGas:   tmConsensusParamUpdates.BlockSize.MaxGas,
			}
		}

		if tmConsensusParamUpdates.Evidence != nil {
			consensusParamUpdates.Evidence = document.EvidenceParams{
				MaxAge: tmConsensusParamUpdates.Evidence.MaxAge,
			}
		}
	}

	res.EndBlock = document.ResponseEndBlock{
		ValidatorUpdates:      validatorUpdates,
		ConsensusParamUpdates: consensusParamUpdates,
		Tags:                  parseTags(result.Results.EndBlock.Tags),
	}

	res.BeginBlock = document.ResponseBeginBlock{
		Tags: parseTags(result.Results.BeginBlock.Tags),
	}

	return res
}

func parseTags(tags []types.TmKVPair) (response []document.KvPair) {
	for _, tag := range tags {
		key := string(tag.Key)
		value := string(tag.Value)
		response = append(response, document.KvPair{Key: key, Value: value})
	}
	return response
}

// parse accounts from coin flow which in block result
// return two kind accounts
// 1. accounts which balance info need updated
// 2. accounts which unbondingDelegation info need updated
func getAccountsFromCoinFlow(endBlockTags []document.KvPair, height int64) ([]string, []string) {
	var (
		accsBalanceNeedUpdated, accsUnbondingDelegationNeedUpdated []string
	)
	balanceAccountExistMap := make(map[string]bool)
	unbondingDelegationAccountExistMap := make(map[string]bool)

	getDistinctAccsBalanceNeedUpdated := func(address string) {
		if strings.HasPrefix(address, bech32AccountAddrPrefix) && !balanceAccountExistMap[address] {
			balanceAccountExistMap[address] = true
			accsBalanceNeedUpdated = append(accsBalanceNeedUpdated, address)
		}
	}
	getDistinctAccsUnbondingDelegationNeedUpdated := func(address string) {
		if strings.HasPrefix(address, bech32AccountAddrPrefix) && !unbondingDelegationAccountExistMap[address] {
			unbondingDelegationAccountExistMap[address] = true
			accsUnbondingDelegationNeedUpdated = append(accsUnbondingDelegationNeedUpdated, address)
		}
	}

	for _, t := range endBlockTags {
		tagKey := string(t.Key)
		tagValue := string(t.Value)

		if assetDetailTriggers[tagKey] || len(tagKey) == triggerTxHashLength {
			values := strings.Split(tagValue, separator)
			if len(values) != 6 {
				logger.Warn("struct of iris coin flow changed in block result, skip parse this block coin flow",
					logger.Int64("height", height), logger.String("tagKey", tagKey))
				continue
			}

			// parse coin flow address from and to, from: value[0], to: value[1]
			from := values[0]
			to := values[1]
			getDistinctAccsBalanceNeedUpdated(from)
			getDistinctAccsBalanceNeedUpdated(to)

			// unbondingDelegation tx complete, need to update account unbondingDelegation info
			if values[3] == unDelegationSubject {
				getDistinctAccsUnbondingDelegationNeedUpdated(from)
				getDistinctAccsUnbondingDelegationNeedUpdated(to)
			}
		}
	}

	return accsBalanceNeedUpdated, accsUnbondingDelegationNeedUpdated
}
