package types

import (
	"fmt"
	"github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/distribution"
	dtags "github.com/irisnet/irishub/app/v1/distribution/tags"
	dtypes "github.com/irisnet/irishub/app/v1/distribution/types"
	"github.com/irisnet/irishub/app/v1/gov"
	"github.com/irisnet/irishub/app/v1/gov/tags"
	"github.com/irisnet/irishub/app/v1/rand"
	"github.com/irisnet/irishub/app/v1/slashing"
	"github.com/irisnet/irishub/app/v1/stake"
	"github.com/irisnet/irishub/app/v2/coinswap"
	"github.com/irisnet/irishub/app/v2/htlc"
	"github.com/irisnet/irishub/modules/guardian"
	stags "github.com/irisnet/irishub/app/v1/stake/tags"
	staketypes "github.com/irisnet/irishub/app/v1/stake/types"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tm "github.com/tendermint/tendermint/types"
	"regexp"
	"strconv"
	"strings"
)

type (
	MsgTransfer      = bank.MsgSend
	MsgBurn          = bank.MsgBurn
	MsgSetMemoRegexp = bank.MsgSetMemoRegexp

	MsgStakeCreate                 = stake.MsgCreateValidator
	MsgStakeEdit                   = stake.MsgEditValidator
	MsgStakeDelegate               = stake.MsgDelegate
	MsgStakeBeginUnbonding         = stake.MsgBeginUnbonding
	MsgBeginRedelegate             = stake.MsgBeginRedelegate
	MsgUnjail                      = slashing.MsgUnjail
	MsgSetWithdrawAddress          = dtypes.MsgSetWithdrawAddress
	MsgWithdrawDelegatorReward     = distribution.MsgWithdrawDelegatorReward
	MsgWithdrawDelegatorRewardsAll = distribution.MsgWithdrawDelegatorRewardsAll
	MsgWithdrawValidatorRewardsAll = distribution.MsgWithdrawValidatorRewardsAll
	StakeValidator                 = stake.Validator
	Delegation                     = stake.Delegation
	UnbondingDelegation            = stake.UnbondingDelegation

	MsgDeposit                       = gov.MsgDeposit
	MsgSubmitProposal                = gov.MsgSubmitProposal
	MsgSubmitSoftwareUpgradeProposal = gov.MsgSubmitSoftwareUpgradeProposal
	MsgSubmitTaxUsageProposal        = gov.MsgSubmitCommunityTaxUsageProposal
	MsgSubmitTokenAdditionProposal   = gov.MsgSubmitTokenAdditionProposal
	MsgVote                          = gov.MsgVote
	Proposal                         = gov.Proposal
	SdkVote                          = gov.Vote

	MsgSwapOrder = coinswap.MsgSwapOrder
	MsgAddLiquidity = coinswap.MsgAddLiquidity
	MsgRemoveLiquidity = coinswap.MsgRemoveLiquidity

	MsgClaimHTLC = htlc.MsgClaimHTLC
	MsgCreateHTLC = htlc.MsgCreateHTLC
	MsgRefundHTLC = htlc.MsgRefundHTLC

	MsgRequestRand = rand.MsgRequestRand

	AssetIssueToken           = asset.MsgIssueToken
	AssetEditToken            = asset.MsgEditToken
	AssetMintToken            = asset.MsgMintToken
	AssetTransferTokenOwner   = asset.MsgTransferTokenOwner
	AssetCreateGateway        = asset.MsgCreateGateway
	AssetEditGateWay          = asset.MsgEditGateway
	AssetTransferGatewayOwner = asset.MsgTransferGatewayOwner

	MsgAddProfiler = guardian.MsgAddProfiler
	MsgAddTrustee = guardian.MsgAddTrustee
	MsgDeleteProfiler = guardian.MsgDeleteProfiler
	MsgDeleteTrustee = guardian.MsgDeleteTrustee

	ResponseDeliverTx = abci.ResponseDeliverTx

	StdTx      = auth.StdTx
	SdkCoins   = types.Coins
	KVPair     = types.KVPair
	AccAddress = types.AccAddress
	ValAddress = types.ValAddress
	Dec        = types.Dec
	Int        = types.Int
	Validator  = tm.Validator
	Tx         = tm.Tx
	Block      = tm.Block
	BlockMeta  = tm.BlockMeta
	HexBytes   = cmn.HexBytes
	TmKVPair   = cmn.KVPair

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
	ValAddressFromBech32 = types.ValAddressFromBech32

	UnmarshalValidator      = staketypes.UnmarshalValidator
	MustUnmarshalValidator  = staketypes.MustUnmarshalValidator
	UnmarshalDelegation     = staketypes.UnmarshalDelegation
	MustUnmarshalDelegation = staketypes.MustUnmarshalDelegation
	MustUnmarshalUBD        = staketypes.MustUnmarshalUBD

	Bech32ifyValPub         = types.Bech32ifyValPub
	Bech32AccountAddrPrefix string
	RegisterCodec           = types.RegisterCodec
	AccAddressFromBech32    = types.AccAddressFromBech32
	BondStatusToString      = types.BondStatusToString

	NewDecFromStr = types.NewDecFromStr

	AddressStoreKey   = auth.AddressStoreKey
	GetAccountDecoder = utils.GetAccountDecoder

	KeyProposal      = gov.KeyProposal
	KeyVotesSubspace = gov.KeyVotesSubspace

	NewHTTP = rpcclient.NewHTTP

	//tags
	TagGovProposalID                   = tags.ProposalID
	TagDistributionReward              = dtags.Reward
	TagStakeActionCompleteRedelegation = stags.ActionCompleteRedelegation
	TagStakeDelegator                  = stags.Delegator
	TagStakeSrcValidator               = stags.SrcValidator
	TagAction                          = types.TagAction

	cdc *codec.Codec
)

// 初始化账户地址前缀
func init() {
	if server.Network == constant.NetworkMainnet {
		types.SetNetworkType(types.Mainnet)
	}
	Bech32AccountAddrPrefix = types.GetConfig().GetBech32AccountAddrPrefix()
	cdc = app.MakeLatestCodec()
}

func GetCodec() *codec.Codec {
	return cdc
}

//
func ParseCoins(coinsStr string) (coins store.Coins) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		coin := ParseCoin(coinStr)
		coins = append(coins, coin)
	}
	return coins
}

func ParseCoin(coinStr string) (coin store.Coin) {
	var (
		reDnm  = `[A-Za-z]{1,}\S*`
		reAmt  = `[0-9]+[.]?[0-9]*`
		reSpc  = `[[:space:]]*`
		reCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnm))
	)

	coinStr = strings.TrimSpace(coinStr)

	matches := reCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		logger.Error("invalid coin expression", logger.Any("coin", coinStr))
		return coin
	}
	denom, amount := matches[2], matches[1]

	amount = getPrecision(amount)
	amt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		logger.Error("Convert str to int failed", logger.Any("amount", amount))
		return coin
	}

	return store.Coin{
		Denom:  denom,
		Amount: amt,
	}
}

func getPrecision(amount string) string {
	length := len(amount)

	if length > 15 {
		nums := strings.Split(amount, ".")
		if len(nums) > 2 {
			return amount
		}

		if len_num0 := len(nums[0]); len_num0 > 15 {
			amount = string([]byte(nums[0])[:15])
			for i := 1; i <= len_num0-15; i++ {
				amount += "0"
			}
		} else {
			leng_append := 16 - len_num0
			amount = nums[0] + "." + string([]byte(nums[1])[:leng_append])
		}
	}
	return amount
}

func BuildFee(fee auth.StdFee) store.Fee {
	return store.Fee{
		Amount: ParseCoins(fee.Amount.String()),
		Gas:    int64(fee.Gas),
	}
}
