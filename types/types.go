package types

import (
	"github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	authcmd "github.com/irisnet/irishub/client/auth/cli"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/gov/tags"
	"github.com/irisnet/irishub/modules/slashing"
	"github.com/irisnet/irishub/modules/stake"
	staketypes "github.com/irisnet/irishub/modules/stake/types"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tm "github.com/tendermint/tendermint/types"
	"strconv"
)

type (
	MsgTransfer = bank.MsgSend

	MsgStakeCreate                 = stake.MsgCreateValidator
	MsgStakeEdit                   = stake.MsgEditValidator
	MsgStakeDelegate               = stake.MsgDelegate
	MsgStakeBeginUnbonding         = stake.MsgBeginUnbonding
	MsgBeginRedelegate             = stake.MsgBeginRedelegate
	MsgUnjail                      = slashing.MsgUnjail
	MsgSetWithdrawAddress          = distribution.MsgSetWithdrawAddress
	MsgWithdrawDelegatorReward     = distribution.MsgWithdrawDelegatorReward
	MsgWithdrawDelegatorRewardsAll = distribution.MsgWithdrawDelegatorRewardsAll
	MsgWithdrawValidatorRewardsAll = distribution.MsgWithdrawValidatorRewardsAll
	StakeValidator                 = stake.Validator
	Delegation                     = stake.Delegation
	UnbondingDelegation            = stake.UnbondingDelegation

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
	ValAddress = types.ValAddress
	Dec        = types.Dec
	Validator  = tm.Validator
	Tx         = tm.Tx
	Block      = tm.Block
	BlockMeta  = tm.BlockMeta
	HexBytes   = cmn.HexBytes

	ABCIQueryOptions = rpcclient.ABCIQueryOptions
	Client           = rpcclient.Client
	HTTP             = rpcclient.HTTP
	ResultStatus     = ctypes.ResultStatus
)

var (
	ValidatorsKey        = stake.ValidatorsKey
	GetValidatorKey      = stake.GetValidatorKey
	GetDelegationKey     = stake.GetDelegationKey
	GetDelegationsKey    = stake.GetDelegationsKey
	GetUBDKey            = stake.GetUBDKey
	GetUBDsKey           = stake.GetUBDsKey
	TagProposalID        = tags.ProposalID
	ValAddressFromBech32 = types.ValAddressFromBech32

	UnmarshalValidator      = staketypes.UnmarshalValidator
	MustUnmarshalValidator  = staketypes.MustUnmarshalValidator
	UnmarshalDelegation     = staketypes.UnmarshalDelegation
	MustUnmarshalDelegation = staketypes.MustUnmarshalDelegation
	MustUnmarshalUBD        = staketypes.MustUnmarshalUBD

	Bech32ifyValPub      = types.Bech32ifyValPub
	RegisterCodec        = types.RegisterCodec
	AccAddressFromBech32 = types.AccAddressFromBech32
	BondStatusToString   = types.BondStatusToString

	NewDecFromStr = types.NewDecFromStr

	AddressStoreKey   = auth.AddressStoreKey
	GetAccountDecoder = authcmd.GetAccountDecoder

	KeyProposal      = gov.KeyProposal
	KeyVotesSubspace = gov.KeyVotesSubspace

	NewHTTP = rpcclient.NewHTTP

	cdc *codec.Codec
)

// 初始化账户地址前缀
func init() {
	config := types.GetConfig()
	config.SetBech32PrefixForAccount(server.Bech32.PrefixAccAddr, server.Bech32.PrefixAccPub)
	config.SetBech32PrefixForValidator(server.Bech32.PrefixValAddr, server.Bech32.PrefixValPub)
	config.SetBech32PrefixForConsensusNode(server.Bech32.PrefixAccAddr, server.Bech32.PrefixConsPub)
	config.Seal()

	cdc = codec.New()

	bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	upgrade.RegisterCodec(cdc)
	distribution.RegisterCodec(cdc)

	types.RegisterCodec(cdc)

	codec.RegisterCrypto(cdc)
}

func GetCodec() *codec.Codec {
	return cdc
}

//
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
