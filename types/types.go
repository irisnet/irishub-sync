package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/stake"
	staketypes "github.com/cosmos/cosmos-sdk/x/stake/types"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tm "github.com/tendermint/tendermint/types"
	"strconv"
)

type (
	MsgTransfer               = bank.MsgSend
	MsgStakeCreate            = stake.MsgCreateValidator
	MsgStakeEdit              = stake.MsgEditValidator
	MsgStakeDelegate          = stake.MsgDelegate
	MsgStakeBeginUnbonding    = stake.MsgBeginUnbonding
	MsgStakeCompleteUnbonding = stake.MsgCompleteUnbonding
	StakeValidator            = stake.Validator
	Delegation                = stake.Delegation
	UnbondingDelegation       = stake.UnbondingDelegation

	MsgDeposit        = gov.MsgDeposit
	MsgSubmitProposal = gov.MsgSubmitProposal
	MsgVote           = gov.MsgVote
	Proposal          = gov.Proposal
	SdkVote           = gov.Vote

	ResponseDeliverTx = abci.ResponseDeliverTx

	StdTx      = auth.StdTx
	SdkCoins   = types.Coins
	KVPair     = types.KVPair
	AccAddress = types.AccAddress
	Validator  = tm.Validator
	Tx         = tm.Tx
	Block      = tm.Block
	HexBytes   = cmn.HexBytes

	Codec            = wire.Codec
	ABCIQueryOptions = rpcclient.ABCIQueryOptions
)

var (
	ValidatorsKey          = stake.ValidatorsKey
	GetValidatorKey        = stake.GetValidatorKey
	GetDelegationKey       = stake.GetDelegationKey
	GetUBDKey              = stake.GetUBDKey
	UnmarshalValidator     = staketypes.UnmarshalValidator
	MustUnmarshalValidator = staketypes.MustUnmarshalValidator
	UnmarshalDelegation    = staketypes.UnmarshalDelegation
	MustUnmarshalUBD       = staketypes.MustUnmarshalUBD

	Bech32ifyValPub      = types.Bech32ifyValPub
	RegisterWire         = types.RegisterWire
	AccAddressFromBech32 = types.AccAddressFromBech32

	AddressStoreKey   = auth.AddressStoreKey
	GetAccountDecoder = authcmd.GetAccountDecoder

	KeyProposal      = gov.KeyProposal
	KeyVotesSubspace = gov.KeyVotesSubspace
)

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
		logger.Error("Can't parse str to float, err is %v\n", logger.String("err", err.Error()))
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
