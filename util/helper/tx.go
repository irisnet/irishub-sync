// package for parse tx struct from binary data

package helper

import (
	"encoding/hex"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	imsg "github.com/irisnet/irishub-sync/msg"
	"github.com/irisnet/irishub-sync/util/constant"
	"strconv"
	"strings"
	"time"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"github.com/irisnet/irishub-sync/msg/nft"
	"github.com/irisnet/irishub-sync/msg/iservice"
	"github.com/irisnet/irishub-sync/msg/oracle"
	"github.com/irisnet/irishub-sync/msg/evidence"
	"github.com/irisnet/irishub-sync/msg/crisis"
)

func ParseTx(txBytes types.Tx, block *types.Block) *document.CommonTx {
	var (
		authTx     types.StdTx
		methodName  = "ParseTx"
		docTx      *document.CommonTx
		gasPrice   float64
		actualFee  store.ActualFee
		signers    []document.Signer
		docTxMsgs  []document.DocTxMsg
	)

	cdc := types.GetCodec()

	err := cdc.UnmarshalBinaryBare(txBytes, &authTx)
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
	//msgData := msgs[0]

	docTx = &document.CommonTx{
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
	for _, msgData := range msgs {
		if len(msgData.GetSigners()) == 0 {
			continue
		}
		if NftTx, ok := nft.HandleTxMsg(msgData, docTx); ok {
			docTx = NftTx
			continue
		}
		if iServiceTx, ok := iservice.HandleTxMsg(msgData, docTx); ok {
			docTx = iServiceTx
			continue
		}
		if OracleTx, ok := oracle.HandleTxMsg(msgData, docTx); ok {
			docTx = OracleTx
			continue
		}
		if EvidenceTx, ok := evidence.HandleTxMsg(msgData, docTx); ok {
			docTx = EvidenceTx
			continue
		}
		if CrisisTx, ok := crisis.HandleTxMsg(msgData, docTx); ok {
			docTx = CrisisTx
			continue
		}

		switch msgData.Type() {
		case new(types.MsgTransfer).Type():
			var msg types.MsgTransfer
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)
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

		case new(types.MsgStakeCreate).Type():
			var msg types.MsgStakeCreate
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

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
		case new(types.MsgStakeEdit).Type():
			var msg types.MsgStakeEdit
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

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

		case new(types.MsgStakeDelegate).Type():
			var msg types.MsgStakeDelegate
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

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

		case new(types.MsgStakeBeginUnbonding).Type():
			var msg types.MsgStakeBeginUnbonding
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

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
		case new(types.MsgBeginRedelegate).Type():
			var msg types.MsgBeginRedelegate
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

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
		case new(types.MsgUnjail).Type():
			var msg types.MsgUnjail
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.ValidatorAddr.String()
			docTx.Type = constant.TxTypeUnjail
			txMsg := imsg.DocTxMsgUnjail{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})
		case new(types.MsgSetWithdrawAddress).Type():
			var msg types.MsgSetWithdrawAddress
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.DelegatorAddress.String()
			docTx.To = msg.WithdrawAddress.String()
			docTx.Type = constant.TxTypeSetWithdrawAddress
			txMsg := imsg.DocTxMsgSetWithdrawAddress{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})
		case new(types.MsgWithdrawDelegatorReward).Type():
			var msg types.MsgWithdrawDelegatorReward
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.DelegatorAddress.String()
			docTx.To = msg.ValidatorAddress.String()
			docTx.Type = constant.TxTypeWithdrawDelegatorReward
			txMsg := imsg.DocTxMsgWithdrawDelegatorReward{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})

		case new(types.MsgFundCommunityPool).Type():
			var msg types.MsgFundCommunityPool
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.Depositor.String()
			docTx.Amount = types.ParseCoins(msg.Amount.String())
			docTx.Type = constant.TxTypeMsgFundCommunityPool
			txMsg := imsg.DocTxMsgFundCommunityPool{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})
		case new(types.MsgWithdrawValidatorCommission).Type():
			var msg types.MsgWithdrawValidatorCommission
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.ValidatorAddress.String()
			docTx.Type = constant.TxTypeMsgWithdrawValidatorCommission
			txMsg := imsg.DocTxMsgWithdrawValidatorCommission{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})

		case new(types.MsgSubmitProposal).Type():
			var msg types.MsgSubmitProposal
			yaml.Unmarshal([]byte(msgData.String()), &msg)

			docTx.Type = constant.TxTypeSubmitProposal
			txMsg := imsg.DocTxMsgSubmitProposal{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})

			//query proposal_id
			proposalId, amount, err := getProposalIdFromEvents(result)
			if err != nil {
				logger.Error("can't get proposal id from tags", logger.String("txHash", docTx.TxHash),
					logger.String("err", err.Error()))
			}
			docTx.ProposalId = proposalId
			docTx.Amount = store.Coins{amount}
			if len(docTx.Signers) > 0 {
				docTx.From = docTx.Signers[0].AddrBech32
			}

		case new(types.MsgDeposit).Type():
			var msg types.MsgDeposit
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

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

		case new(types.MsgVote).Type():
			var msg types.MsgVote
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

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
		case new(types.MsgRequestRandom).Type():
			var msg types.MsgRequestRandom
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.Consumer.String()
			docTx.Amount = []store.Coin{}
			docTx.Type = constant.TxTypeRequestRand
			txMsg := imsg.DocTxMsgRequestRand{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})
		case new(types.MsgIssueToken).Type():
			var msg types.MsgIssueToken
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.Owner.String()
			docTx.Type = constant.TxTypeAssetIssueToken
			txMsg := imsg.DocTxMsgIssueToken{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})

		case new(types.MsgEditToken).Type():
			var msg types.MsgEditToken
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.Owner.String()
			docTx.Type = constant.TxTypeAssetEditToken
			txMsg := imsg.DocTxMsgEditToken{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})

		case new(types.MsgMintToken).Type():
			var msg types.MsgMintToken
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.Owner.String()
			docTx.To = msg.To.String()
			docTx.Type = constant.TxTypeAssetMintToken
			txMsg := imsg.DocTxMsgMintToken{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})

		case new(types.MsgTransferTokenOwner).Type():
			var msg types.MsgTransferTokenOwner
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.SrcOwner.String()
			docTx.To = msg.DstOwner.String()
			docTx.Type = constant.TxTypeAssetTransferTokenOwner
			txMsg := imsg.DocTxMsgTransferTokenOwner{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})

		case new(types.MsgAddProfiler).Type():
			var msg types.MsgAddProfiler
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.AddGuardian.AddedBy.String()
			docTx.To = msg.AddGuardian.Address.String()
			docTx.Type = constant.TxTypeAddProfiler
			txMsg := imsg.DocTxMsgAddProfiler{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})

		case new(types.MsgAddTrustee).Type():
			var msg types.MsgAddTrustee
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.AddGuardian.AddedBy.String()
			docTx.To = msg.AddGuardian.Address.String()
			docTx.Type = constant.TxTypeAddTrustee
			txMsg := imsg.DocTxMsgAddTrustee{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})

		case new(types.MsgDeleteTrustee).Type():
			var msg types.MsgDeleteTrustee
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.DeleteGuardian.DeletedBy.String()
			docTx.To = msg.DeleteGuardian.Address.String()
			docTx.Type = constant.TxTypeDeleteTrustee
			txMsg := imsg.DocTxMsgDeleteTrustee{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})

		case new(types.MsgDeleteProfiler).Type():
			var msg types.MsgDeleteProfiler
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.DeleteGuardian.DeletedBy.String()
			docTx.To = msg.DeleteGuardian.Address.String()
			docTx.Type = constant.TxTypeDeleteProfiler
			txMsg := imsg.DocTxMsgDeleteProfiler{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})

		case new(types.MsgCreateHTLC).Type():
			var msg types.MsgCreateHTLC
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

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
		case new(types.MsgClaimHTLC).Type():
			var msg types.MsgClaimHTLC
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.Sender.String()
			docTx.To = ""
			docTx.Type = constant.TxTypeClaimHTLC
			txMsg := imsg.DocTxMsgClaimHTLC{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})
		case new(types.MsgRefundHTLC).Type():
			var msg types.MsgRefundHTLC
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

			docTx.From = msg.Sender.String()
			docTx.To = ""
			docTx.Type = constant.TxTypeRefundHTLC
			txMsg := imsg.DocTxMsgRefundHTLC{}
			txMsg.BuildMsg(msg)
			docTx.Msgs = append(docTxMsgs, document.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})
		case new(types.MsgAddLiquidity).Type():
			var msg types.MsgAddLiquidity
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

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
		case new(types.MsgRemoveLiquidity).Type():
			var msg types.MsgRemoveLiquidity
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

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
		case new(types.MsgSwapOrder).Type():
			var msg types.MsgSwapOrder
			data, _ := json.Marshal(msgData)
			json.Unmarshal(data, &msg)

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

		default:
			logger.Warn("unknown msg type")
		}
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

//func getProposerFromEvents(result types.ResponseDeliverTx) (string) {
//	for _, val := range result.GetEvents() {
//		if val.Type != "message" {
//			continue
//		}
//		for _, attr := range val.Attributes {
//			if string(attr.Key) == "sender" {
//				return string(attr.Value)
//			}
//		}
//	}
//	return ""
//}

// get proposalId from tags
func getProposalIdFromEvents(result types.ResponseDeliverTx) (uint64, store.Coin, error) {
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
	var proposalId uint64
	var amount store.Coin
	for _, val := range result.GetEvents() {
		if val.Type != types.EventTypeProposalDeposit {
			continue
		}
		for _, attr := range val.Attributes {
			if string(attr.Key) == types.EventGovProposalID {
				if id, err := strconv.ParseInt(string(attr.Value), 10, 0); err == nil {
					proposalId = uint64(id)
				}
			}
			if string(attr.Key) == "amount" && string(attr.Value) != "" {
				value := string(attr.Value)
				amount = types.ParseCoin(value)
			}
		}
	}

	return proposalId, amount, nil
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
