package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"strconv"
)

type (
	MsgTransfer               = bank.MsgSend
	MsgStakeCreate            = stake.MsgCreateValidator
	MsgStakeEdit              = stake.MsgEditValidator
	MsgStakeDelegate          = stake.MsgDelegate
	MsgStakeBeginUnbonding    = stake.MsgBeginUnbonding
	MsgStakeCompleteUnbonding = stake.MsgCompleteUnbonding
	MsgDeposit                = gov.MsgDeposit
	MsgSubmitProposal         = gov.MsgSubmitProposal
	MsgVote                   = gov.MsgVote
)

type Msg interface {
	Type() string
	String() string
}

func BuildCoins(coins types.Coins) store.Coins {
	var (
		localCoins store.Coins
	)

	if len(coins) > 0 {
		for _, coin := range coins {
			localCoins = append(localCoins, BuildCoin(coin))
		}
	}

	return localCoins
}

func BuildCoin(coin types.Coin) store.Coin {
	amount, err := strconv.ParseFloat(coin.Amount.String(), 64)
	if err != nil {
		logger.Error.Printf("Can't parse str to float, err is %v\n", err)
	}
	return store.Coin{
		Denom:  coin.Denom,
		Amount: amount,
	}
}

func BuildFee(fee auth.StdFee) store.Fee {
	return store.Fee{
		Amount: BuildCoins(fee.Amount),
		Gas:    fee.Gas,
	}
}
