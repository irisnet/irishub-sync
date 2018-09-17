// package for parse tx struct from binary data

package helper

import (
	"encoding/hex"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	"strings"
)

func ParseTx(cdc *itypes.Codec, txBytes itypes.Tx, block *itypes.Block) document.CommonTx {
	var (
		authTx     itypes.StdTx
		methodName = "ParseTx"
		docTx      document.CommonTx
		gasPrice   float64
		actualFee  store.ActualFee
	)

	err := cdc.UnmarshalBinary(txBytes, &authTx)
	if err != nil {
		logger.Error.Println(err)
		return docTx
	}

	height := block.Height
	time := block.Time
	txHash := BuildHex(txBytes.Hash())
	fee := itypes.BuildFee(authTx.Fee)
	memo := authTx.Memo

	// get tx status, gasUsed, gasPrice and actualFee from tx result
	status, result, err := QueryTxResult(txBytes.Hash())
	if err != nil {
		logger.Error.Printf("%v: can't get txResult, err is %v\n", methodName, err)
	}
	log := result.Log
	gasUsed := result.GasUsed
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
		logger.Warning.Printf("%v: can't get msgs\n", methodName)
		return docTx
	}
	msg := msgs[0]

	docTx = document.CommonTx{
		Height:    height,
		Time:      time,
		TxHash:    txHash,
		Fee:       fee,
		Memo:      memo,
		Status:    status,
		Log:       log,
		GasUsed:   gasUsed,
		GasPrice:  gasPrice,
		ActualFee: actualFee,
	}

	switch msg.(type) {
	case itypes.MsgTransfer:
		msg := msg.(itypes.MsgTransfer)

		docTx.From = msg.Inputs[0].Address.String()
		docTx.To = msg.Outputs[0].Address.String()
		docTx.Amount = itypes.BuildCoins(msg.Inputs[0].Coins)
		docTx.Type = constant.TxTypeTransfer
		return docTx
	case itypes.MsgStakeCreate:
		msg := msg.(itypes.MsgStakeCreate)

		docTx.From = msg.ValidatorAddr.String()
		docTx.To = ""
		docTx.Amount = []store.Coin{itypes.BuildCoin(msg.Delegation)}
		docTx.Type = constant.TxTypeStakeCreateValidator

		// struct of createValidator
		valDes := document.ValDescription{
			Moniker:  msg.Moniker,
			Identity: msg.Identity,
			Website:  msg.Website,
			Details:  msg.Details,
		}
		pubKey, err := itypes.Bech32ifyValPub(msg.PubKey)
		if err != nil {
			logger.Error.Printf("%v: Can't get pubKey, txHash is %v\n",
				methodName, txHash)
			pubKey = ""
		}
		docTx.StakeCreateValidator = document.StakeCreateValidator{
			PubKey:      pubKey,
			Description: valDes,
		}

		return docTx
	case itypes.MsgStakeEdit:
		msg := msg.(itypes.MsgStakeEdit)

		docTx.From = msg.ValidatorAddr.String()
		docTx.To = ""
		docTx.Amount = []store.Coin{}
		docTx.Type = constant.TxTypeStakeEditValidator

		// struct of editValidator
		valDes := document.ValDescription{
			Moniker:  msg.Moniker,
			Identity: msg.Identity,
			Website:  msg.Website,
			Details:  msg.Details,
		}
		docTx.StakeEditValidator = document.StakeEditValidator{
			Description: valDes,
		}

		return docTx
	case itypes.MsgStakeDelegate:
		msg := msg.(itypes.MsgStakeDelegate)

		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()
		docTx.Amount = []store.Coin{itypes.BuildCoin(msg.Delegation)}
		docTx.Type = constant.TxTypeStakeDelegate

		return docTx
	case itypes.MsgStakeBeginUnbonding:
		msg := msg.(itypes.MsgStakeBeginUnbonding)

		shares, _ := msg.SharesAmount.Float64()
		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()

		coin := store.Coin{
			Amount: shares,
		}
		docTx.Amount = []store.Coin{coin}
		docTx.Type = constant.TxTypeStakeBeginUnbonding
		return docTx
	case itypes.MsgStakeCompleteUnbonding:
		msg := msg.(itypes.MsgStakeCompleteUnbonding)

		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()
		docTx.Amount = nil
		docTx.Type = constant.TxTypeStakeCompleteUnbonding
		return docTx
	case itypes.MsgSubmitProposal:
		msg := msg.(itypes.MsgSubmitProposal)

		docTx.From = msg.Proposer.String()
		docTx.To = ""
		docTx.Amount = itypes.BuildCoins(msg.InitialDeposit)
		docTx.Type = constant.TxTypeSubmitProposal
		docTx.Msg = itypes.NewSubmitProposal(msg)

		//query proposal_id
		for _, tag := range result.Tags {
			key := string(tag.Key)
			if key == "proposalId" {
				var proposalId int64
				cdc.MustUnmarshalBinaryBare(tag.Value, &proposalId)
				docTx.ProposalId = proposalId
			}
		}
		return docTx
	case itypes.MsgDeposit:
		msg := msg.(itypes.MsgDeposit)

		docTx.From = msg.Depositer.String()
		docTx.Amount = itypes.BuildCoins(msg.Amount)
		docTx.Type = constant.TxTypeDeposit
		docTx.Msg = itypes.NewDeposit(msg)
		docTx.ProposalId = msg.ProposalID
		return docTx
	case itypes.MsgVote:
		msg := msg.(itypes.MsgVote)

		docTx.From = msg.Voter.String()
		docTx.Amount = []store.Coin{}
		docTx.Type = constant.TxTypeVote
		docTx.Msg = itypes.NewVote(msg)
		docTx.ProposalId = msg.ProposalID
		return docTx

	default:
		logger.Info.Println("unknown msg type")
	}

	return docTx
}
func BuildHex(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

// get tx status and log by query txHash
func QueryTxResult(txHash []byte) (string, itypes.ResponseDeliverTx, error) {
	var resDeliverTx itypes.ResponseDeliverTx
	status := constant.TxStatusSuccess

	client := GetClient()
	defer client.Release()

	res, err := client.Client.Tx(txHash, false)
	if err != nil {
		return "unknown", resDeliverTx, err
	}
	result := res.TxResult
	if result.Code != 0 {
		status = constant.TxStatusFail
	}

	return status, result, nil
}
