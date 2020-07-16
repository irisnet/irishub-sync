package types

import (
	"fmt"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub/app"
	asset "github.com/irismod/token/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	rand "github.com/irisnet/irishub/modules/random/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stake "github.com/cosmos/cosmos-sdk/x/staking/types"
	staketypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	coinswap "github.com/irismod/coinswap/types"
	htlc "github.com/irismod/htlc/types"
	"github.com/cosmos/cosmos-sdk/codec"
	guardian "github.com/irisnet/irishub/modules/guardian/types"
	"github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/bytes"
	cmnk "github.com/tendermint/tendermint/crypto/merkle"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpcclienthttp "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tm "github.com/tendermint/tendermint/types"
	"regexp"
	"strconv"
	"strings"
)

type (
	MsgTransfer = bank.MsgSend

	MsgStakeCreate = stake.MsgCreateValidator
	MsgStakeEdit = stake.MsgEditValidator
	MsgStakeDelegate = stake.MsgDelegate
	MsgStakeBeginUnbonding = stake.MsgUndelegate
	MsgBeginRedelegate = stake.MsgBeginRedelegate
	MsgUnjail = slashing.MsgUnjail
	MsgSetWithdrawAddress = dtypes.MsgSetWithdrawAddress
	MsgWithdrawDelegatorReward = distribution.MsgWithdrawDelegatorReward
	MsgFundCommunityPool = distribution.MsgFundCommunityPool
	MsgWithdrawValidatorCommission = distribution.MsgWithdrawValidatorCommission
	StakeValidator = stake.Validator
	Delegation = stake.Delegation
	UnbondingDelegation = stake.UnbondingDelegation

	MsgDeposit = gov.MsgDeposit
	MsgSubmitProposal = gov.MsgSubmitProposal
	MsgVote = gov.MsgVote
	Proposal = gov.Proposal
	SdkVote = gov.Vote

	MsgSwapOrder = coinswap.MsgSwapOrder
	MsgAddLiquidity = coinswap.MsgAddLiquidity
	MsgRemoveLiquidity = coinswap.MsgRemoveLiquidity

	MsgClaimHTLC = htlc.MsgClaimHTLC
	MsgCreateHTLC = htlc.MsgCreateHTLC
	MsgRefundHTLC = htlc.MsgRefundHTLC

	MsgRequestRandom = rand.MsgRequestRandom

	AssetIssueToken = asset.MsgIssueToken
	AssetEditToken = asset.MsgEditToken
	AssetMintToken = asset.MsgMintToken
	AssetTransferTokenOwner = asset.MsgTransferTokenOwner

	MsgAddProfiler = guardian.MsgAddProfiler
	MsgAddTrustee = guardian.MsgAddTrustee
	MsgDeleteProfiler = guardian.MsgDeleteProfiler
	MsgDeleteTrustee = guardian.MsgDeleteTrustee

	ResponseDeliverTx = abci.ResponseDeliverTx

	StdTx = auth.StdTx
	SdkCoins = types.Coins
	KVPair = types.KVPair
	AccAddress = types.AccAddress
	ValAddress = types.ValAddress
	Dec = types.Dec
	Int = types.Int
	Validator = tm.Validator
	Tx = tm.Tx
	Block = tm.Block
	BlockID = tm.BlockID
	//BlockMeta = tm.BlockMeta
	HexBytes = cmn.HexBytes
	TmKVPair = cmnk.KVPair

	ABCIQueryOptions = rpcclient.ABCIQueryOptions
	Client = rpcclient.Client
	ResultStatus = ctypes.ResultStatus
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

	//Bech32ifyValPub         = types.Bech32ifyValPub
	Bech32AccountAddrPrefix string
	RegisterCodec           = types.RegisterCodec
	AccAddressFromBech32    = types.AccAddressFromBech32
	//BondStatusToString      = types.BondStatusToString

	NewDecFromStr = types.NewDecFromStr

	//AddressStoreKey   = auth.AddressStoreKey
	//StoreName         = auth.StoreKey
	//GetAccountDecoder = auth.NewAccountRetriever

	KeyProposal      = gov.ProposalKey
	KeyVotesSubspace = gov.VotesKey
	//
	NewHTTP = rpcclienthttp.New

	//tags
	EventGovProposalID = gov.AttributeKeyProposalID
	//TagDistributionReward              = dtags.Reward
	//TagStakeActionCompleteRedelegation = stags.ActionCompleteRedelegation
	//TagStakeDelegator                  = stags.Delegator
	//TagStakeSrcValidator               = stags.SrcValidator
	//TagAction                          = types.TagAction

	cdc *codec.Codec
)

// 初始化账户地址前缀
func init() {
	Bech32AccountAddrPrefix = types.GetConfig().GetBech32AccountAddrPrefix()
	_, cdc = app.MakeCodecs()
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
