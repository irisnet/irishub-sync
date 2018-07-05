// package for parse tx struct from binary data

package helper

import (
	"github.com/irisnet/irishub-sync/model/store/document"

	"github.com/tendermint/tendermint/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/stake"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/model/store"
	"strings"
	"encoding/hex"
	"github.com/irisnet/irishub-sync/util/constant"
	"strconv"
)

type (
	msgBankSend = bank.MsgSend
	msgStakeCreate = stake.MsgCreateValidator
	msgStakeEdit = stake.MsgEditValidator
	msgStakeDelegate = stake.MsgDelegate
	msgStakeUnbond = stake.MsgUnbond
)

func ParseTx(cdc *wire.Codec, txBytes types.Tx, block *types.Block) store.Docs {
	var (
		authTx auth.StdTx
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
	status := constant.TxStatusSuccess

	switch authTx.GetMsg().(type) {
	case msgBankSend:
		msg := authTx.Msg.(msgBankSend)
		docTx := document.CommonTx {
			Height: height,
			Time: time,
			TxHash: txHash,
			Fee: fee,
			Status: status,
		}
		docTx.From = msg.Inputs[0].Address.String()
		docTx.To = msg.Outputs[0].Address.String()
		docTx.Amount = BuildCoins(msg.Inputs[0].Coins)
		docTx.Type = constant.TxTypeBank
		return docTx
	case msgStakeCreate:
		msg := authTx.Msg.(msgStakeCreate)
		stakeTx := document.StakeTx{
			Height: height,
			Time: time,
			TxHash: txHash,
			Fee: fee,
			Status: status,
		}
		stakeTx.ValidatorAddr = msg.ValidatorAddr.String()
		stakeTx.PubKey = BuildHex(msg.PubKey.Bytes())
		stakeTx.Amount = buildCoin(msg.Bond)

		description := document.Description{
			Moniker: msg.Moniker,
			Identity: msg.Identity,
			Website: msg.Website,
			Details: msg.Details,
		}

		docTx := document.StakeTxDeclareCandidacy{
			StakeTx: stakeTx,
			Description: description,
		}
		docTx.Type = constant.TxTypeStakeCreate
		return docTx
	case msgStakeEdit:
		msg := authTx.Msg.(msgStakeEdit)
		stakeTx := document.StakeTx{
			Height: height,
			Time: time,
			TxHash: txHash,
			Fee: fee,
			Status: status,
		}
		stakeTx.ValidatorAddr = msg.ValidatorAddr.String()

		description := document.Description{
			Moniker: msg.Moniker,
			Identity: msg.Identity,
			Website: msg.Website,
			Details: msg.Details,
		}

		docTx := document.StakeTxEditCandidacy{
			StakeTx: stakeTx,
			Description: description,
		}
		docTx.Type = constant.TxTypeStakeEdit
		return docTx
	case msgStakeDelegate:
		msg := authTx.Msg.(msgStakeDelegate)
		docTx := document.StakeTx{
			Height: height,
			Time: time,
			TxHash: txHash,
			Fee: fee,
			Status: status,
		}
		docTx.DelegatorAddr = msg.DelegatorAddr.String()
		docTx.ValidatorAddr = msg.ValidatorAddr.String()
		docTx.Amount = buildCoin(msg.Bond)
		docTx.Type = constant.TxTypeStakeDelegate
		return docTx
	case msgStakeUnbond:
		msg := authTx.Msg.(msgStakeUnbond)
		shares, err := strconv.Atoi(msg.Shares)
		if err != nil {
			logger.Error.Println(err)
		}
		docTx := document.StakeTx{
			Height: height,
			Time: time,
			TxHash: txHash,
			Fee: fee,
			Status: status,
		}
		docTx.DelegatorAddr = msg.DelegatorAddr.String()
		docTx.ValidatorAddr = msg.ValidatorAddr.String()
		docTx.Amount = store.Coin{
			Amount: int64(shares),
		}
		docTx.Type = constant.TxTypeStakeUnbond
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
		Denom: coin.Denom,
		Amount: coin.Amount,
	}
}

func buildFee(fee auth.StdFee) store.Fee {
	return store.Fee{
		Amount: BuildCoins(fee.Amount),
		Gas:    fee.Gas,
	}
}

func BuildHex(bytes []byte) string  {
	return strings.ToUpper(hex.EncodeToString(bytes))
}
