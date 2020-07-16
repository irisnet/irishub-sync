// package for parse tx struct from binary data

package helper

import (
	"encoding/hex"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	imsg "github.com/irisnet/irishub-sync/types/msg"
	"github.com/irisnet/irishub-sync/util/constant"
	"strconv"
	"strings"
	"time"
)

func ParseTx(txBytes types.Tx, block *types.Block) document.CommonTx {
	var (
		authTx     types.StdTx
		methodName = "ParseTx"
		docTx      document.CommonTx
		gasPrice   float64
		actualFee  store.ActualFee
		signers    []document.Signer
		docTxMsgs  []document.DocTxMsg
	)

	cdc := types.GetCodec()

	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &authTx)
	if err != nil {
		logger.Error(err.Error())
		return docTx
	}

	height := block.Height
	blockTime := block.Time
	txHash := BuildHex(txBytes.Hash())
	fee := types.BuildFee(authTx.Fee)
	memo := authTx.Memo

	// get tx signers
	if len(authTx.Signatures) > 0 {
		for _, signature := range authTx.GetSigners() {
			signer := document.Signer{}
			signer.AddrHex = signature.String()
			if addrBech32, err := ConvertAccountAddrFromHexToBech32(signature.Bytes()); err != nil {
				logger.Error("convert account addr from hex to bech32 fail",
					logger.String("addrHex", signature.String()), logger.String("err", err.Error()))
			} else {
				signer.AddrBech32 = addrBech32
			}
			signers = append(signers, signer)
		}
	}

	// get tx status, gasUsed, gasPrice and actualFee from tx result
	status, result, err := QueryTxResult(txBytes.Hash())
	if err != nil {
		logger.Error("get txResult err", logger.String("method", methodName), logger.String("err", err.Error()))
	}
	log := result.Log
	gasUsed := Min(result.GasUsed, fee.Gas)
	if len(fee.Amount) > 0 {
		gasPrice = fee.Amount[0].Amount / float64(fee.Gas)
		actualFee = store.ActualFee{
			Denom:  fee.Amount[0].Denom,
			Amount: float64(gasUsed) * gasPrice,
		}
	} else {
		gasPrice = 0
		actualFee = store.ActualFee{}
	}

	msgs := authTx.GetMsgs()
	if len(msgs) <= 0 {
		logger.Error("can't get msgs", logger.String("method", methodName))
		return docTx
	}
	msg := msgs[0]

	docTx = document.CommonTx{
		Height:    height,
		Time:      blockTime,
		TxHash:    txHash,
		Fee:       fee,
		Memo:      memo,
		Status:    status,
		Code:      result.Code,
		Log:       log,
		GasUsed:   gasUsed,
		GasWanted: result.GasUsed,
		GasPrice:  gasPrice,
		ActualFee: actualFee,
		Events:    parseEvents(result),
		Signers:   signers,
	}

	switch msg.(type) {
	case types.MsgTransfer:
		msg := msg.(types.MsgTransfer)

		docTx.From = msg.FromAddress.String()
		docTx.To = msg.ToAddress.String()
		docTx.Amount = types.ParseCoins(msg.Amount.String())
		docTx.Type = constant.TxTypeTransfer
		txMsg := imsg.DocTxMsgSend{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx

	case types.MsgStakeCreate:
		msg := msg.(types.MsgStakeCreate)

		docTx.From = msg.DelegatorAddress.String()
		docTx.To = msg.ValidatorAddress.String()
		docTx.Amount = []store.Coin{types.ParseCoin(msg.Value.String())}
		docTx.Type = constant.TxTypeStakeCreateValidator
		txMsg := imsg.DocTxMsgStakeCreate{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.MsgStakeEdit:
		msg := msg.(types.MsgStakeEdit)

		docTx.From = msg.ValidatorAddress.String()
		docTx.To = ""
		docTx.Amount = []store.Coin{}
		docTx.Type = constant.TxTypeStakeEditValidator
		txMsg := imsg.DocTxMsgStakeEdit{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.MsgStakeDelegate:
		msg := msg.(types.MsgStakeDelegate)

		docTx.From = msg.DelegatorAddress.String()
		docTx.To = msg.ValidatorAddress.String()
		docTx.Amount = []store.Coin{types.ParseCoin(msg.Amount.String())}
		docTx.Type = constant.TxTypeStakeDelegate
		txMsg := imsg.DocTxMsgDelegate{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx
	case types.MsgStakeBeginUnbonding:
		msg := msg.(types.MsgStakeBeginUnbonding)

		shares := ParseFloat(msg.Amount.String())
		docTx.From = msg.DelegatorAddress.String()
		docTx.To = msg.ValidatorAddress.String()

		coin := store.Coin{
			Amount: shares,
		}
		docTx.Amount = []store.Coin{coin}
		docTx.Type = constant.TxTypeStakeBeginUnbonding
		txMsg := imsg.DocTxMsgBeginUnbonding{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.MsgBeginRedelegate:
		msg := msg.(types.MsgBeginRedelegate)

		shares := ParseFloat(msg.Amount.String())
		docTx.From = msg.ValidatorSrcAddress.String()
		docTx.To = msg.ValidatorDstAddress.String()
		coin := store.Coin{
			Amount: shares,
		}
		docTx.Amount = []store.Coin{coin}
		docTx.Type = constant.TxTypeBeginRedelegate
		txMsg := imsg.DocTxMsgBeginRedelegate{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.MsgUnjail:
		msg := msg.(types.MsgUnjail)

		docTx.From = msg.ValidatorAddr.String()
		docTx.Type = constant.TxTypeUnjail
		txMsg := imsg.DocTxMsgUnjail{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
	case types.MsgSetWithdrawAddress:
		msg := msg.(types.MsgSetWithdrawAddress)

		docTx.From = msg.DelegatorAddress.String()
		docTx.To = msg.WithdrawAddress.String()
		docTx.Type = constant.TxTypeSetWithdrawAddress
		txMsg := imsg.DocTxMsgSetWithdrawAddress{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
	case types.MsgWithdrawDelegatorReward:
		msg := msg.(types.MsgWithdrawDelegatorReward)

		docTx.From = msg.DelegatorAddress.String()
		docTx.To = msg.ValidatorAddress.String()
		docTx.Type = constant.TxTypeWithdrawDelegatorReward
		txMsg := imsg.DocTxMsgWithdrawDelegatorReward{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

	case types.MsgFundCommunityPool:
		msg := msg.(types.MsgFundCommunityPool)

		docTx.From = msg.Depositor.String()
		docTx.Amount = types.ParseCoins(msg.Amount.String())
		docTx.Type = constant.TxTypeMsgFundCommunityPool
		txMsg := imsg.DocTxMsgFundCommunityPool{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
	case types.MsgWithdrawValidatorCommission:
		msg := msg.(types.MsgWithdrawValidatorCommission)

		docTx.From = msg.ValidatorAddress.String()
		docTx.Type = constant.TxTypeMsgWithdrawValidatorCommission
		txMsg := imsg.DocTxMsgWithdrawValidatorCommission{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

	case types.MsgSubmitProposal:
		msg := msg.(types.MsgSubmitProposal)

		docTx.From = msg.Proposer.String()
		docTx.To = ""
		docTx.Amount = types.ParseCoins(msg.InitialDeposit.String())
		docTx.Type = constant.TxTypeSubmitProposal
		txMsg := imsg.DocTxMsgSubmitProposal{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		//query proposal_id
		proposalId, err := getProposalIdFromTags(result)
		if err != nil {
			logger.Error("can't get proposal id from tags", logger.String("txHash", docTx.TxHash),
				logger.String("err", err.Error()))
		}
		docTx.ProposalId = proposalId

		return docTx
		//case types.MsgSubmitSoftwareUpgradeProposal:
		//	msg := msg.(types.MsgSubmitSoftwareUpgradeProposal)
		//
		//	docTx.From = msg.Proposer.String()
		//	docTx.To = ""
		//	docTx.Amount = types.ParseCoins(msg.InitialDeposit.String())
		//	docTx.Type = constant.TxTypeSubmitProposal
		//	txMsg := imsg.DocTxMsgSubmitSoftwareUpgradeProposal{}
		//	txMsg.BuildMsg(msg)
		//	docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
		//		Type: txMsg.Type(),
		//		Msg:  &txMsg,
		//	})
		//
		//	//query proposal_id
		//	proposalId, err := getProposalIdFromTags(result.Tags)
		//	if err != nil {
		//		logger.Error("can't get proposal id from tags", logger.String("txHash", docTx.TxHash),
		//			logger.String("err", err.Error()))
		//	}
		//	docTx.ProposalId = proposalId
		//
		//	return docTx
		//case types.MsgSubmitTaxUsageProposal:
		//	msg := msg.(types.MsgSubmitTaxUsageProposal)
		//
		//	docTx.From = msg.Proposer.String()
		//	docTx.To = ""
		//	docTx.Amount = types.ParseCoins(msg.InitialDeposit.String())
		//	docTx.Type = constant.TxTypeSubmitProposal
		//	txMsg := imsg.DocTxMsgSubmitCommunityTaxUsageProposal{}
		//	txMsg.BuildMsg(msg)
		//	docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
		//		Type: txMsg.Type(),
		//		Msg:  &txMsg,
		//	})
		//
		//	//query proposal_id
		//	proposalId, err := getProposalIdFromTags(result.Tags)
		//	if err != nil {
		//		logger.Error("can't get proposal id from tags", logger.String("txHash", docTx.TxHash),
		//			logger.String("err", err.Error()))
		//	}
		//	docTx.ProposalId = proposalId
		//	return docTx
		//case types.MsgSubmitTokenAdditionProposal:
		//	msg := msg.(types.MsgSubmitTokenAdditionProposal)
		//
		//	docTx.From = msg.Proposer.String()
		//	docTx.To = ""
		//	docTx.Amount = types.ParseCoins(msg.InitialDeposit.String())
		//	docTx.Type = constant.TxTypeSubmitProposal
		//	txMsg := imsg.DocTxMsgSubmitTokenAdditionProposal{}
		//	txMsg.BuildMsg(msg)
		//	docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
		//		Type: txMsg.Type(),
		//		Msg:  &txMsg,
		//	})
		//	//query proposal_id
		//	proposalId, err := getProposalIdFromTags(result.Tags)
		//	if err != nil {
		//		logger.Error("can't get proposal id from tags", logger.String("txHash", docTx.TxHash),
		//			logger.String("err", err.Error()))
		//	}
		//	docTx.ProposalId = proposalId
		//	return docTx
	case types.MsgDeposit:
		msg := msg.(types.MsgDeposit)

		docTx.From = msg.Depositor.String()
		docTx.Amount = types.ParseCoins(msg.Amount.String())
		docTx.Type = constant.TxTypeDeposit
		docTx.ProposalId = msg.ProposalID
		txMsg := imsg.DocTxMsgDeposit{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.MsgVote:
		msg := msg.(types.MsgVote)

		docTx.From = msg.Voter.String()
		docTx.Amount = []store.Coin{}
		docTx.Type = constant.TxTypeVote
		docTx.ProposalId = msg.ProposalID
		txMsg := imsg.DocTxMsgVote{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.MsgRequestRandom:
		msg := msg.(types.MsgRequestRandom)

		docTx.From = msg.Consumer.String()
		docTx.Amount = []store.Coin{}
		docTx.Type = constant.TxTypeRequestRand
		txMsg := imsg.DocTxMsgRequestRand{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.AssetIssueToken:
		msg := msg.(types.AssetIssueToken)

		docTx.From = msg.Owner.String()
		docTx.Type = constant.TxTypeAssetIssueToken
		txMsg := imsg.DocTxMsgIssueToken{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx
	case types.AssetEditToken:
		msg := msg.(types.AssetEditToken)

		docTx.From = msg.Owner.String()
		docTx.Type = constant.TxTypeAssetEditToken
		txMsg := imsg.DocTxMsgEditToken{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx
	case types.AssetMintToken:
		msg := msg.(types.AssetMintToken)

		docTx.From = msg.Owner.String()
		docTx.To = msg.To.String()
		docTx.Type = constant.TxTypeAssetMintToken
		txMsg := imsg.DocTxMsgMintToken{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx
	case types.AssetTransferTokenOwner:
		msg := msg.(types.AssetTransferTokenOwner)

		docTx.From = msg.SrcOwner.String()
		docTx.To = msg.DstOwner.String()
		docTx.Type = constant.TxTypeAssetTransferTokenOwner
		txMsg := imsg.DocTxMsgTransferTokenOwner{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx
		//case types.AssetCreateGateway:
		//	msg := msg.(types.AssetCreateGateway)
		//
		//	docTx.From = msg.Owner.String()
		//	docTx.Type = constant.TxTypeAssetCreateGateway
		//	txMsg := imsg.DocTxMsgCreateGateway{}
		//	txMsg.BuildMsg(msg)
		//	docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
		//		Type: txMsg.Type(),
		//		Msg:  &txMsg,
		//	})
		//
		//	return docTx
		//case types.AssetEditGateWay:
		//	msg := msg.(types.AssetEditGateWay)
		//
		//	docTx.From = msg.Owner.String()
		//	docTx.Type = constant.TxTypeAssetEditGateway
		//	txMsg := imsg.DocTxMsgEditGateway{}
		//	txMsg.BuildMsg(msg)
		//	docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
		//		Type: txMsg.Type(),
		//		Msg:  &txMsg,
		//	})
		//
		//	return docTx
		//case types.AssetTransferGatewayOwner:
		//	msg := msg.(types.AssetTransferGatewayOwner)
		//
		//	docTx.From = msg.Owner.String()
		//	docTx.To = msg.To.String()
		//	docTx.Type = constant.TxTypeAssetTransferGatewayOwner
		//	txMsg := imsg.DocTxMsgTransferGatewayOwner{}
		//	txMsg.BuildMsg(msg)
		//	docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
		//		Type: txMsg.Type(),
		//		Msg:  &txMsg,
		//	})
		//	return docTx

	case types.MsgAddProfiler:
		msg := msg.(types.MsgAddProfiler)

		docTx.From = msg.AddGuardian.AddedBy.String()
		docTx.To = msg.AddGuardian.Address.String()
		docTx.Type = constant.TxTypeAddProfiler
		txMsg := imsg.DocTxMsgAddProfiler{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx

	case types.MsgAddTrustee:
		msg := msg.(types.MsgAddTrustee)

		docTx.From = msg.AddGuardian.AddedBy.String()
		docTx.To = msg.AddGuardian.Address.String()
		docTx.Type = constant.TxTypeAddTrustee
		txMsg := imsg.DocTxMsgAddTrustee{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx

	case types.MsgDeleteTrustee:
		msg := msg.(types.MsgDeleteTrustee)

		docTx.From = msg.DeleteGuardian.DeletedBy.String()
		docTx.To = msg.DeleteGuardian.Address.String()
		docTx.Type = constant.TxTypeDeleteTrustee
		txMsg := imsg.DocTxMsgDeleteTrustee{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx

	case types.MsgDeleteProfiler:
		msg := msg.(types.MsgDeleteProfiler)

		docTx.From = msg.DeleteGuardian.DeletedBy.String()
		docTx.To = msg.DeleteGuardian.Address.String()
		docTx.Type = constant.TxTypeDeleteProfiler
		txMsg := imsg.DocTxMsgDeleteProfiler{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx

	case types.MsgCreateHTLC:
		msg := msg.(types.MsgCreateHTLC)

		docTx.From = msg.Sender.String()
		docTx.To = msg.To.String()
		docTx.Amount = types.ParseCoins(msg.Amount.String())
		docTx.Type = constant.TxTypeCreateHTLC
		txMsg := imsg.DocTxMsgCreateHTLC{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.MsgClaimHTLC:
		msg := msg.(types.MsgClaimHTLC)

		docTx.From = msg.Sender.String()
		docTx.To = ""
		docTx.Type = constant.TxTypeClaimHTLC
		txMsg := imsg.DocTxMsgClaimHTLC{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.MsgRefundHTLC:
		msg := msg.(types.MsgRefundHTLC)

		docTx.From = msg.Sender.String()
		docTx.To = ""
		docTx.Type = constant.TxTypeRefundHTLC
		txMsg := imsg.DocTxMsgRefundHTLC{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.MsgAddLiquidity:
		msg := msg.(types.MsgAddLiquidity)

		docTx.From = msg.Sender.String()
		docTx.To = ""
		docTx.Amount = types.ParseCoins(msg.MaxToken.String())
		docTx.Type = constant.TxTypeAddLiquidity
		txMsg := imsg.DocTxMsgAddLiquidity{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.MsgRemoveLiquidity:
		msg := msg.(types.MsgRemoveLiquidity)

		docTx.From = msg.Sender.String()
		docTx.To = ""
		docTx.Amount = types.ParseCoins(msg.WithdrawLiquidity.String())
		docTx.Type = constant.TxTypeRemoveLiquidity
		txMsg := imsg.DocTxMsgRemoveLiquidity{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	case types.MsgSwapOrder:
		msg := msg.(types.MsgSwapOrder)

		docTx.From = msg.Input.Address.String()
		docTx.To = msg.Output.Address.String()
		docTx.Amount = types.ParseCoins(msg.Input.Coin.String())
		docTx.Type = constant.TxTypeSwapOrder
		txMsg := imsg.DocTxMsgSwapOrder{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx

	default:
		logger.Warn("unknown msg type")
	}

	return docTx
}

func parseEvents(result types.ResponseDeliverTx) []document.Event {

	var events []document.Event
	for _, val := range result.GetEvents() {
		one := document.Event{
			Type: val.Type,
		}
		one.Attributes = make(map[string]string, len(val.Attributes))
		for _, attr := range val.Attributes {
			one.Attributes[string(attr.Key)] = string(attr.Value)
		}
		events = append(events, one)
	}

	return events
}

// get proposalId from tags
func getProposalIdFromTags(result types.ResponseDeliverTx) (uint64, error) {
	//query proposal_id
	//for _, tag := range tags {
	//	key := string(tag.Key)
	//	if key == types.EventGovProposalID {
	//		if proposalId, err := strconv.ParseInt(string(tag.Value), 10, 0); err != nil {
	//			return 0, err
	//		} else {
	//			return uint64(proposalId), nil
	//		}
	//	}
	//}
	for _, val := range result.GetEvents() {
		for key, attr := range val.Attributes {
			if string(key) == types.EventGovProposalID {
				if proposalId, err := strconv.ParseInt(string(attr.Value), 10, 0); err != nil {
					return 0, err
				} else {
					return uint64(proposalId), nil
				}
			}
		}
	}
	return 0, nil
}

func BuildHex(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

// get tx status and log by query txHash
func QueryTxResult(txHash []byte) (string, types.ResponseDeliverTx, error) {
	var resDeliverTx types.ResponseDeliverTx
	status := document.TxStatusSuccess

	client := GetClient()
	defer client.Release()

	res, err := client.Tx(txHash, false)
	if err != nil {
		// try again
		time.Sleep(time.Duration(1) * time.Second)
		if res, err := client.Tx(txHash, false); err != nil {
			return "unknown", resDeliverTx, err
		} else {
			resDeliverTx = res.TxResult
		}
	} else {
		resDeliverTx = res.TxResult
	}

	if resDeliverTx.Code != 0 {
		status = document.TxStatusFail
	}

	return status, resDeliverTx, nil
}
