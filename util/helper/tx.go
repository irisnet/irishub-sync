// package for parse tx struct from binary data

package helper

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/irisnet/iris-sync-server/module/logger"
	"github.com/irisnet/iris-sync-server/util/constant"

	sdk "github.com/cosmos/cosmos-sdk"
	"github.com/tendermint/go-wire/data"
	"github.com/tendermint/tendermint/types"

	// 需要将 cosmos-sdk module 中的 txInner 注册到内存
	// 中，才能解析 tx 结构
	_ "github.com/cosmos/cosmos-sdk/modules/auth"
	_ "github.com/cosmos/cosmos-sdk/modules/base"
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"github.com/cosmos/cosmos-sdk/modules/nonce"
	"github.com/irisnet/iris-sync-server/module/stake"
	"github.com/irisnet/iris-sync-server/model/store/document"
)

// parse tx struct from binary data
func ParseTx(txByte types.Tx) (string, interface{}) {
	txb, err := sdk.LoadTx(txByte)
	if err != nil {
		logger.Error.Println(err)
	}
	txl, ok := txb.Unwrap().(sdk.TxLayer)

	var (
		txi       sdk.Tx
		coinTx    document.CoinTx
		stakeTx   document.StakeTx
		StakeTxDeclareCandidacy document.StakeTxDeclareCandidacy
		nonceAddr data.Bytes
	)

	for ok {
		txi = txl.Next()

		switch txi.Unwrap().(type) {

		case coin.SendTx:
			ctx, _ := txi.Unwrap().(coin.SendTx)
			coinTx.From = fmt.Sprintf("%s", ctx.Inputs[0].Address.Address)
			coinTx.To = fmt.Sprintf("%s", ctx.Outputs[0].Address.Address)
			coinTx.Amount = ctx.Inputs[0].Coins
			coinTx.TxHash = strings.ToUpper(hex.EncodeToString(txByte.Hash()))
			return constant.TxTypeCoin, coinTx
		case nonce.Tx:
			ctx, _ := txi.Unwrap().(nonce.Tx)
			nonceAddr = ctx.Signers[0].Address
			fmt.Println(nonceAddr)
			break
		case stake.TxUnbond, stake.TxDelegate, stake.TxDeclareCandidacy:
			kind, _ := txi.GetKind()
			stakeTx.From = fmt.Sprintf("%s", nonceAddr)
			stakeTx.Type = strings.Replace(kind, constant.TxTypeStake + "/", "", -1)
			stakeTx.TxHash = strings.ToUpper(hex.EncodeToString(txByte.Hash()))

			switch kind {
			case stake.TypeTxDeclareCandidacy:
				ctx, _ := txi.Unwrap().(stake.TxDeclareCandidacy)
				stakeTx.Amount.Denom = ctx.BondUpdate.Bond.Denom
				stakeTx.Amount.Amount = ctx.BondUpdate.Bond.Amount
				stakeTx.PubKey = fmt.Sprintf("%s", ctx.PubKey.KeyString())

				description := document.Description{
					Moniker: ctx.Description.Moniker,
					Identity: ctx.Description.Identity,
					Website: ctx.Description.Website,
					Details: ctx.Description.Details,
				}

				StakeTxDeclareCandidacy.StakeTx = stakeTx
				StakeTxDeclareCandidacy.Description = description

				return kind, StakeTxDeclareCandidacy
			case stake.TypeTxEditCandidacy:
				// TODO：record edit candidacy tx if necessary
				// ctx, _ := txi.Unwrap().(stake.TxEditCandidacy)
				break
			case stake.TypeTxDelegate:
				ctx, _ := txi.Unwrap().(stake.TxDelegate)
				stakeTx.Amount.Denom = ctx.Bond.Denom
				stakeTx.Amount.Amount = ctx.Bond.Amount
				stakeTx.PubKey = fmt.Sprintf("%s", ctx.PubKey.KeyString())
				break
			case stake.TypeTxUnbond:
				ctx, _ := txi.Unwrap().(stake.TxUnbond)
				stakeTx.Amount.Amount = int64(ctx.Shares)
				stakeTx.PubKey = fmt.Sprintf("%s", ctx.PubKey.KeyString())
				break
			}
			return kind, stakeTx
		default:
			// logger.Info.Printf("unsupported tx type, %+v\n", txi.Unwrap())
		}

		txl, ok = txi.Unwrap().(sdk.TxLayer)
	}
	return "", nil
}
