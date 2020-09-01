package handler

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/helper"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/irisnet/irishub-sync/store"
)


func ParseBlock(meta *types.BlockID, block *types.Block, validators []*types.Validator) document.Block {
	//cdc := types.GetCodec()

	hexFunc := func(bytes []byte) string {
		return helper.BuildHex(bytes)
	}

	docBlock := document.Block{
		Height:          block.Header.Height,
		Hash:            hexFunc(meta.Hash),
		Time:            block.Header.Time,
		NumTxs:          int64(len(block.Data.Txs)),
		ProposalAddress: block.Header.ProposerAddress.String(),
	}

	//lastBlockId := document.BlockID{
	//	Hash: hexFunc(meta.Header.LastBlockID.Hash),
	//	PartsHeader: document.PartSetHeader{
	//		Total: meta.Header.LastBlockID.PartsHeader.Total,
	//		Hash:  hexFunc(meta.Header.LastBlockID.PartsHeader.Hash),
	//	},
	//}

	// blockMeta
	blockMeta := document.BlockMeta{
		//BlockID: document.BlockID{
		//	Hash: hexFunc(meta.BlockID.Hash),
		//	PartsHeader: document.PartSetHeader{
		//		Total: meta.BlockID.PartsHeader.Total,
		//		Hash:  hexFunc(meta.BlockID.PartsHeader.Hash),
		//	},
		//},
		Header: document.Header{
			//ChainID:         meta.Header.ChainID,
			//Height:          meta.Header.Height,
			//Time:            meta.Header.Time,
			//NumTxs:          meta.Header.NumTxs,
			//LastBlockID:     lastBlockId,
			//TotalTxs: int64(block.NumTxs),
			//LastCommitHash:  hexFunc(meta.Header.LastCommitHash),
			//DataHash:        hexFunc(meta.Header.DataHash),
			//ValidatorsHash:  hexFunc(meta.Header.ValidatorsHash),
			//ConsensusHash:   hexFunc(meta.Header.ConsensusHash),
			//AppHash:         hexFunc(meta.Header.AppHash),
			//LastResultsHash: hexFunc(meta.Header.LastResultsHash),
			//EvidenceHash:    hexFunc(meta.Header.EvidenceHash),
		},
	}

	// block
	var (
		preCommits []document.Vote
	)

	if len(block.LastCommit.Signatures) > 0 {
		for idx, v := range block.LastCommit.Signatures {
			if v.Signature != nil {
				//var sig document.Signature
				//out, _ := cdc.MarshalJSON(v.Signature)
				//json.Unmarshal(out, &sig)
				vote := block.LastCommit.GetVote(int32(idx))
				preCommit := document.Vote{
					ValidatorAddress: vote.ValidatorAddress.String(),
					ValidatorIndex:   vote.ValidatorIndex,
					Height:           vote.Height,
					//Round:            vote.Round,
					Timestamp: vote.Timestamp,
					//Type:             byte(vote.Type),
					//BlockID:          vote.BlockID,
					//Signature:        sig,
				}
				preCommits = append(preCommits, preCommit)
			}
		}
	}

	blockContent := document.BlockContent{
		LastCommit: document.Commit{
			//BlockID:    lastBlockId,
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
	parseBlockResult(docBlock.Height)

	return docBlock
}

func parseBlockResult(height int64) {
	client := helper.GetClient()
	defer client.Release()

	result, err := client.BlockResults(&height)
	if err != nil {
		// try again
		var err2 error
		client2 := helper.GetClient()
		result, err2 = client2.BlockResults(&height)
		client2.Release()
		if err2 != nil {
			logger.Error("parse block result fail", logger.Int64("block", height),
				logger.String("err", err.Error()))
			return
		}
	}

	if proposalId, ok := IsContainVotingEndEvent(parseEvents(result.EndBlockEvents)); ok {
		if proposal, err := document.QueryProposal(proposalId); err == nil {
			proposal.VotingEndHeight = height
			store.SaveOrUpdate(proposal)
		} else {
			logger.Error("QueryProposal fail", logger.Int64("block", height),
				logger.String("err", err.Error()))
		}
	}

	return
}

func parseEvents(events []abcitypes.Event) (response []document.Event) {
	for _, event := range events {
		one := document.Event{Type: event.Type}
		for _, v := range event.Attributes {
			one.Attributes = append(one.Attributes, document.Attribute{Key: string(v.Key), Value: string(v.Value)})
		}
		response = append(response, one)
	}
	return response
}

//// parse accounts from coin flow which in block result
//// return two kind accounts
//// 1. accounts which balance info need updated
//// 2. accounts which unbondingDelegation info need updated
//func getAccountsFromCoinFlow(endBlockTags []document.KvPair, height int64) ([]string, []string) {
//	var (
//		accsBalanceNeedUpdated, accsUnbondingDelegationNeedUpdated []string
//	)
//	balanceAccountExistMap := make(map[string]bool)
//	unbondingDelegationAccountExistMap := make(map[string]bool)
//
//	getDistinctAccsBalanceNeedUpdated := func(address string) {
//		if strings.HasPrefix(address, bech32AccountAddrPrefix) && !balanceAccountExistMap[address] {
//			balanceAccountExistMap[address] = true
//			accsBalanceNeedUpdated = append(accsBalanceNeedUpdated, address)
//		}
//	}
//	getDistinctAccsUnbondingDelegationNeedUpdated := func(address string) {
//		if strings.HasPrefix(address, bech32AccountAddrPrefix) && !unbondingDelegationAccountExistMap[address] {
//			unbondingDelegationAccountExistMap[address] = true
//			accsUnbondingDelegationNeedUpdated = append(accsUnbondingDelegationNeedUpdated, address)
//		}
//	}
//
//	for _, t := range endBlockTags {
//		tagKey := string(t.Key)
//		tagValue := string(t.Value)
//
//		if assetDetailTriggers[tagKey] || len(tagKey) == triggerTxHashLength {
//			values := strings.Split(tagValue, separator)
//			if len(values) != 6 {
//				logger.Warn("struct of iris coin flow changed in block result, skip parse this block coin flow",
//					logger.Int64("height", height), logger.String("tagKey", tagKey))
//				continue
//			}
//
//			// parse coin flow address from and to, from: value[0], to: value[1]
//			from := values[0]
//			to := values[1]
//			getDistinctAccsBalanceNeedUpdated(from)
//			getDistinctAccsBalanceNeedUpdated(to)
//
//			// unbondingDelegation tx complete, need to update account unbondingDelegation info
//			if values[3] == unDelegationSubject {
//				getDistinctAccsUnbondingDelegationNeedUpdated(from)
//				getDistinctAccsUnbondingDelegationNeedUpdated(to)
//			}
//		}
//	}
//
//	return accsBalanceNeedUpdated, accsUnbondingDelegationNeedUpdated
//}
