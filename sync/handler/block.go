package handler

import (
	"github.com/irisnet/irishub-sync/util/helper"
	"github.com/irisnet/irishub-sync/model/store/document"
	"github.com/tendermint/tendermint/types"
	"github.com/irisnet/irishub-sync/model/store"
	"github.com/irisnet/irishub-sync/module/codec"
	"encoding/json"
	"github.com/irisnet/irishub-sync/module/logger"
)

func SaveBlock(meta *types.BlockMeta, block *types.Block)  {

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
			Hash: hexFunc(meta.Header.LastBlockID.PartsHeader.Hash),
		},
	}

	blockMeta := document.BlockMeta{
		BlockID: document.BlockID{
			Hash: hexFunc(meta.BlockID.Hash),
			PartsHeader: document.PartSetHeader{
				Total: meta.BlockID.PartsHeader.Total,
				Hash: hexFunc(meta.BlockID.PartsHeader.Hash),
			},
		},
		Header: document.Header{
			ChainID: meta.Header.ChainID,
			Height: meta.Header.Height,
			Time: meta.Header.Time,
			NumTxs: meta.Header.NumTxs,
			LastBlockID: lastBlockId,
			TotalTxs: meta.Header.TotalTxs,
			LastCommitHash: hexFunc(meta.Header.LastCommitHash),
			DataHash: hexFunc(meta.Header.DataHash),
			ValidatorsHash: hexFunc(meta.Header.ValidatorsHash),
			ConsensusHash: hexFunc(meta.Header.ConsensusHash),
			AppHash: hexFunc(meta.Header.AppHash),
			LastResultsHash: hexFunc(meta.Header.LastResultsHash),
			EvidenceHash: hexFunc(meta.Header.EvidenceHash),
		},
	}

	var (
		preCommits []document.Vote
	)

	if len(block.LastCommit.Precommits) > 0 {
		for _, v := range block.LastCommit.Precommits {
			var sig document.Signature
			out, _ := codec.Cdc.MarshalJSON(v.Signature)
			json.Unmarshal(out, &sig)
			preCommit := document.Vote{
				ValidatorAddress: v.ValidatorAddress.String(),
				ValidatorIndex: v.ValidatorIndex,
				Height: v.Height,
				Round: v.Round,
				Timestamp: v.Timestamp,
				Type: v.Type,
				BlockID: lastBlockId,
				Signature: sig,
			}
			preCommits = append(preCommits, preCommit)
		}
	}

	blockContent := document.BlockContent{
		LastCommit: document.Commit{
			BlockID: lastBlockId,
			Precommits: preCommits,
		},
	}

	docBlock.Meta = blockMeta
	docBlock.Block = blockContent

	err := store.Save(docBlock)
	if err != nil {
		logger.Error.Println(err)
	}
}
