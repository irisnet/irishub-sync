package distribution

import (
	types "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store"
)

// msg struct for changing the withdraw address for a delegator (or validator self-delegation)
type DocTxMsgSetWithdrawAddress struct {
	DelegatorAddr string `bson:"delegator_addr"`
	WithdrawAddr  string `bson:"withdraw_addr"`
}

func (doctx *DocTxMsgSetWithdrawAddress) Type() string {
	return constant.TxTypeSetWithdrawAddress
}

func (doctx *DocTxMsgSetWithdrawAddress) BuildMsg(txMsg interface{}) {
	msg := txMsg.(types.MsgSetWithdrawAddress)
	doctx.DelegatorAddr = msg.DelegatorAddress.String()
	doctx.WithdrawAddr = msg.WithdrawAddress.String()
}

// msg struct for delegation withdraw from a single validator
type DocTxMsgWithdrawDelegatorReward struct {
	DelegatorAddr string `bson:"delegator_addr"`
	ValidatorAddr string `bson:"validator_addr"`
}

func (doctx *DocTxMsgWithdrawDelegatorReward) Type() string {
	return constant.TxTypeWithdrawDelegatorReward
}

func (doctx *DocTxMsgWithdrawDelegatorReward) BuildMsg(txMsg interface{}) {
	msg := txMsg.(types.MsgWithdrawDelegatorReward)
	doctx.DelegatorAddr = msg.DelegatorAddress.String()
	doctx.ValidatorAddr = msg.ValidatorAddress.String()
}

// msg struct for delegation withdraw for all of the delegator's delegations
type DocTxMsgFundCommunityPool struct {
	Amount    store.Coins `bson:"amount"`
	Depositor string      `bson:"depositor"`
}

func (doctx *DocTxMsgFundCommunityPool) Type() string {
	return constant.TxTypeMsgFundCommunityPool
}

func (doctx *DocTxMsgFundCommunityPool) BuildMsg(txMsg interface{}) {
	msg := txMsg.(types.MsgFundCommunityPool)
	doctx.Depositor = msg.Depositor.String()
	doctx.Amount = types.ParseCoins(msg.Amount.String())
}

// msg struct for validator withdraw
type DocTxMsgWithdrawValidatorCommission struct {
	ValidatorAddr string `bson:"validator_addr"`
}

func (doctx *DocTxMsgWithdrawValidatorCommission) Type() string {
	return constant.TxTypeMsgWithdrawValidatorCommission
}

func (doctx *DocTxMsgWithdrawValidatorCommission) BuildMsg(txMsg interface{}) {
	msg := txMsg.(types.MsgWithdrawValidatorCommission)
	doctx.ValidatorAddr = msg.ValidatorAddress.String()
}
