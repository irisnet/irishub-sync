// package for parse tx struct from binary data

package helper

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
	"strings"
)

type (
	msgTransfer               = bank.MsgSend
	msgStakeCreate            = stake.MsgCreateValidator
	msgStakeEdit              = stake.MsgEditValidator
	msgStakeDelegate          = stake.MsgDelegate
	msgStakeBeginUnbonding    = stake.MsgBeginUnbonding
	msgStakeCompleteUnbonding = stake.MsgCompleteUnbonding
)

func ParseTx(cdc *wire.Codec, txBytes types.Tx, block *types.Block) document.CommonTx {
	var (
		authTx     auth.StdTx
		methodName = "ParseTx"
		docTx      document.CommonTx
		gasPrice   float64
	)

	err := cdc.UnmarshalBinary(txBytes, &authTx)
	if err != nil {
		logger.Error.Println(err)
		return docTx
	}

	height := block.Height
	time := block.Time
	txHash := BuildHex(txBytes.Hash())
	fee := buildFee(authTx.Fee)
	memo := authTx.Memo

	// get tx status, gasUsed, gasPrice from tx result
	status, result, err := getTxResult(txBytes.Hash())
	if err != nil {
		logger.Error.Printf("%v: can't get txResult, err is %v\n", methodName, err)
	}
	log := result.Log
	gasUsed := result.GasUsed
	if len(fee.Amount) > 0 {
		gasPrice = fee.Amount[0].Amount / float64(fee.Gas)
	} else {
		gasPrice = 0
	}

	msgs := authTx.GetMsgs()
	if len(msgs) <= 0 {
		logger.Warning.Printf("%v: can't get msgs\n", methodName)
		return docTx
	}
	msg := msgs[0]

	switch msg.(type) {
	case msgTransfer:
		msg := msg.(msgTransfer)
		docTx = document.CommonTx{
			Height:   height,
			Time:     time,
			TxHash:   txHash,
			Fee:      fee,
			Memo:     memo,
			Status:   status,
			Log:      log,
			GasUsed:  gasUsed,
			GasPrice: gasPrice,
		}
		docTx.From = msg.Inputs[0].Address.String()
		docTx.To = msg.Outputs[0].Address.String()
		docTx.Amount = BuildCoins(msg.Inputs[0].Coins)
		docTx.Type = constant.TxTypeTransfer
		return docTx
	case msgStakeCreate:
		msg := msg.(msgStakeCreate)
		docTx = document.CommonTx{
			Height:   height,
			Time:     time,
			TxHash:   txHash,
			Fee:      fee,
			Memo:     memo,
			Status:   status,
			Log:      log,
			GasUsed:  gasUsed,
			GasPrice: gasPrice,
		}
		docTx.From = msg.ValidatorAddr.String()
		docTx.To = ""
		docTx.Amount = []store.Coin{buildCoin(msg.Delegation)}
		docTx.Type = constant.TxTypeStakeCreateValidator

		// struct of createValidator
		valDes := document.ValDescription{
			Moniker:  msg.Moniker,
			Identity: msg.Identity,
			Website:  msg.Website,
			Details:  msg.Details,
		}
		pubKey, err := sdk.Bech32ifyValPub(msg.PubKey)
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
	case msgStakeEdit:
		msg := msg.(msgStakeEdit)
		docTx = document.CommonTx{
			Height:   height,
			Time:     time,
			TxHash:   txHash,
			Fee:      fee,
			Memo:     memo,
			Status:   status,
			Log:      log,
			GasUsed:  gasUsed,
			GasPrice: gasPrice,
		}
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
	case msgStakeDelegate:
		msg := msg.(msgStakeDelegate)
		docTx = document.CommonTx{
			Height:   height,
			Time:     time,
			TxHash:   txHash,
			Fee:      fee,
			Memo:     memo,
			Status:   status,
			Log:      log,
			GasUsed:  gasUsed,
			GasPrice: gasPrice,
		}
		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()
		docTx.Amount = []store.Coin{buildCoin(msg.Delegation)}
		docTx.Type = constant.TxTypeStakeDelegate

		return docTx
	case msgStakeBeginUnbonding:
		msg := msg.(msgStakeBeginUnbonding)
		shares, _ := msg.SharesAmount.Float64()

		docTx = document.CommonTx{
			Height:   height,
			Time:     time,
			TxHash:   txHash,
			Fee:      fee,
			Memo:     memo,
			Status:   status,
			Log:      log,
			GasUsed:  gasUsed,
			GasPrice: gasPrice,
		}
		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()

		coin := store.Coin{
			Amount: shares,
		}
		docTx.Amount = []store.Coin{coin}
		docTx.Type = constant.TxTypeStakeBeginUnbonding
		return docTx
	case msgStakeCompleteUnbonding:
		msg := msg.(msgStakeCompleteUnbonding)

		docTx := document.CommonTx{
			Height:   height,
			Time:     time,
			TxHash:   txHash,
			Fee:      fee,
			Memo:     memo,
			Status:   status,
			Log:      log,
			GasUsed:  gasUsed,
			GasPrice: gasPrice,
		}
		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()
		docTx.Amount = nil
		docTx.Type = constant.TxTypeStakeCompleteUnbonding
		return docTx
	default:
		logger.Info.Println("unknown msg type")
	}

	return docTx
}

func BuildCoins(coins sdktypes.Coins) store.Coins {
	var (
		localCoins store.Coins
	)

	if len(coins) > 0 {
		for _, coin := range coins {
			localCoins = append(localCoins, buildCoin(coin))
		}
	}

	return localCoins
}

func buildCoin(coin sdktypes.Coin) store.Coin {
	amount, err := ParseStrToFloat(coin.Amount.String())
	if err != nil {
		logger.Error.Printf("Can't parse str to float, err is %v\n", err)
	}
	return store.Coin{
		Denom:  coin.Denom,
		Amount: amount,
	}
}

func buildFee(fee auth.StdFee) store.Fee {
	return store.Fee{
		Amount: BuildCoins(fee.Amount),
		Gas:    fee.Gas,
	}
}

func BuildHex(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

// get tx status and log by query txHash
func getTxResult(txHash []byte) (string, abci.ResponseDeliverTx, error) {
	var resDeliverTx abci.ResponseDeliverTx
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
