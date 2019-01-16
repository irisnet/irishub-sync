package handler

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/helper"
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
