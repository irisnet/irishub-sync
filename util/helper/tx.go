package helper

import (
	"encoding/hex"
	"fmt"
	"strings"

	conf "github.com/irisnet/iris-sync-server/conf/server"
	"github.com/irisnet/iris-sync-server/model/store/collection"
	"github.com/irisnet/iris-sync-server/module/logger"

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
		coinTx    collection.CoinTx
		stakeTx   collection.StakeTx
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
			return "coin", coinTx
		case nonce.Tx:
			ctx, _ := txi.Unwrap().(nonce.Tx)
			nonceAddr = ctx.Signers[0].Address
			fmt.Println(nonceAddr)
			break
		case stake.TxUnbond, stake.TxDelegate, stake.TxDeclareCandidacy:
			kind, _ := txi.GetKind()
			stakeTx.From = fmt.Sprintf("%s", nonceAddr)
			stakeTx.Type = strings.Replace(kind, "stake/", "", -1)
			stakeTx.TxHash = strings.ToUpper(hex.EncodeToString(txByte.Hash()))

			switch kind {
			case stake.TypeTxDeclareCandidacy:
				ctx, _ := txi.Unwrap().(stake.TxDeclareCandidacy)
				stakeTx.Amount.Denom = ctx.BondUpdate.Bond.Denom
				stakeTx.Amount.Amount = ctx.BondUpdate.Bond.Amount
				stakeTx.PubKey = fmt.Sprintf("%s", ctx.PubKey.KeyString())
				break
			case stake.TypeTxEditCandidacy:
				// TODO：record edit candidacy tx if necessary
				//ctx, _ := txi.Unwrap().(stake.TxEditCandidacy)
				break
			case stake.TypeTxDelegate:
				ctx, _ := txi.Unwrap().(stake.TxDelegate)
				stakeTx.Amount.Denom = ctx.Bond.Denom
				stakeTx.Amount.Amount = ctx.Bond.Amount
				stakeTx.PubKey = fmt.Sprintf("%s", ctx.PubKey.KeyString())
				break
			case stake.TypeTxUnbond:
				ctx, _ := txi.Unwrap().(stake.TxUnbond)
				stakeTx.Amount.Denom = conf.Token
				stakeTx.Amount.Amount = int64(ctx.Shares)
				stakeTx.PubKey = fmt.Sprintf("%s", ctx.PubKey.KeyString())
				break
			}
			return "stake", stakeTx
		default:
			logger.Info.Println("unsupported tx type")
		}

		txl, ok = txi.Unwrap().(sdk.TxLayer)
	}
	return "", nil
}
