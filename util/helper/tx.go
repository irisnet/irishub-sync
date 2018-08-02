// package for parse tx struct from binary data

package helper

import (
	"github.com/irisnet/irishub-sync/store/document"

	"encoding/hex"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/constant"
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

func ParseTx(cdc *wire.Codec, txBytes types.Tx, block *types.Block) store.Docs {
	var (
		authTx     auth.StdTx
		methodName = "ParseTx"
	)

	err := cdc.UnmarshalBinary(txBytes, &authTx)
	if err != nil {
		logger.Error.Println(err)
		return nil
	}

	height := block.Height
	time := block.Time
	txHash := BuildHex(txBytes.Hash())
	fee := buildFee(authTx.Fee)
	memo := authTx.Memo
	status, log, err := getTxResult(txBytes.Hash())
	if err != nil {
		logger.Error.Printf("%v: can't get txResult, err is %v\n", methodName, err)
	}

	msgs := authTx.GetMsgs()
	if len(msgs) <= 0 {
		logger.Warning.Printf("%v: can't get msgs\n", methodName)
		return nil
	}
	msg := msgs[0]

	switch msg.(type) {
	case msgTransfer:
		msg := msg.(msgTransfer)
		docTx := document.CommonTx{
			Height: height,
			Time:   time,
			TxHash: txHash,
			Fee:    fee,
			Memo:   memo,
			Status: status,
			Log:    log,
		}
		docTx.From = msg.Inputs[0].Address.String()
		docTx.To = msg.Outputs[0].Address.String()
		docTx.Amount = BuildCoins(msg.Inputs[0].Coins)
		docTx.Type = constant.TxTypeTransfer
		return docTx
	case msgStakeCreate:
		msg := msg.(msgStakeCreate)
		stakeTx := document.StakeTx{
			Height: height,
			Time:   time,
			TxHash: txHash,
			Fee:    fee,
			Memo:   memo,
			Status: status,
			Log:    log,
		}
		stakeTx.ValidatorAddr = msg.ValidatorAddr.String()
		stakeTx.DelegatorAddr = msg.DelegatorAddr.String()
		stakeTx.PubKey = BuildHex(msg.PubKey.Bytes())
		stakeTx.Amount = buildCoin(msg.Delegation)

		description := document.Description{
			Moniker:  msg.Moniker,
			Identity: msg.Identity,
			Website:  msg.Website,
			Details:  msg.Details,
		}

		docTx := document.StakeTxDeclareCandidacy{
			StakeTx:     stakeTx,
			Description: description,
		}
		docTx.Type = constant.TxTypeStakeCreateValidator
		return docTx
	case msgStakeEdit:
		msg := msg.(msgStakeEdit)
		stakeTx := document.StakeTx{
			Height: height,
			Time:   time,
			TxHash: txHash,
			Fee:    fee,
			Memo:   memo,
			Status: status,
			Log:    log,
		}
		stakeTx.ValidatorAddr = msg.ValidatorAddr.String()

		description := document.Description{
			Moniker:  msg.Moniker,
			Identity: msg.Identity,
			Website:  msg.Website,
			Details:  msg.Details,
		}

		docTx := document.StakeTxEditCandidacy{
			StakeTx:     stakeTx,
			Description: description,
		}
		docTx.Type = constant.TxTypeStakeEditValidator
		return docTx
	case msgStakeDelegate:
		msg := msg.(msgStakeDelegate)
		docTx := document.StakeTx{
			Height: height,
			Time:   time,
			TxHash: txHash,
			Fee:    fee,
			Memo:   memo,
			Status: status,
			Log:    log,
		}
		docTx.DelegatorAddr = msg.DelegatorAddr.String()
		docTx.ValidatorAddr = msg.ValidatorAddr.String()
		docTx.Amount = buildCoin(msg.Delegation)
		docTx.Type = constant.TxTypeStakeDelegate
		return docTx
	case msgStakeBeginUnbonding:
		msg := msg.(msgStakeBeginUnbonding)
		shares, _ := msg.SharesAmount.Float64()

		docTx := document.StakeTx{
			Height: height,
			Time:   time,
			TxHash: txHash,
			Fee:    fee,
			Memo:   memo,
			Status: status,
			Log:    log,
		}
		docTx.DelegatorAddr = msg.DelegatorAddr.String()
		docTx.ValidatorAddr = msg.ValidatorAddr.String()
		docTx.Amount = store.Coin{
			Amount: shares,
		}
		docTx.Type = constant.TxTypeStakeBeginUnbonding
		return docTx
	case msgStakeCompleteUnbonding:
		msg := msg.(msgStakeCompleteUnbonding)

		docTx := document.StakeTx{
			Height: height,
			Time:   time,
			TxHash: txHash,
			Fee:    fee,
			Memo:   memo,
			Status: status,
			Log:    log,
		}
		docTx.DelegatorAddr = msg.DelegatorAddr.String()
		docTx.ValidatorAddr = msg.ValidatorAddr.String()
		docTx.Type = constant.TxTypeStakeCompleteUnbonding
		return docTx
	default:
		logger.Info.Println("unknown msg type")
	}

	return nil
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
	return store.Coin{
		Denom:  coin.Denom,
		Amount: float64(coin.Amount.Int64()),
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
func getTxResult(txHash []byte) (string, string, error) {

	status := constant.TxStatusSuccess

	client := GetClient()
	defer client.Release()

	res, err := client.Client.Tx(txHash, false)
	if err != nil {
		return "unknown", "", err
	}
	result := res.TxResult
	if result.Code != 0 {
		status = constant.TxStatusFail
	}

	return status, result.Log, nil
}
