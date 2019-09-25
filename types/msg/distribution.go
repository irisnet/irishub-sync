package msg

import (
	"github.com/irisnet/irishub-sync/util/constant"
	itypes "github.com/irisnet/irishub-sync/types"
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
	msg := txMsg.(itypes.MsgSetWithdrawAddress)
	doctx.DelegatorAddr = msg.DelegatorAddr.String()
	doctx.WithdrawAddr = msg.WithdrawAddr.String()
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
	msg := txMsg.(itypes.MsgWithdrawDelegatorReward)
	doctx.DelegatorAddr = msg.DelegatorAddr.String()
	doctx.ValidatorAddr = msg.ValidatorAddr.String()
}

// msg struct for delegation withdraw for all of the delegator's delegations
type DocTxMsgWithdrawDelegatorRewardsAll struct {
	DelegatorAddr string `bson:"delegator_addr"`
}

func (doctx *DocTxMsgWithdrawDelegatorRewardsAll) Type() string {
	return constant.TxTypeWithdrawDelegatorRewardsAll
}

func (doctx *DocTxMsgWithdrawDelegatorRewardsAll) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgWithdrawDelegatorRewardsAll)
	doctx.DelegatorAddr = msg.DelegatorAddr.String()
}

// msg struct for validator withdraw
type DocTxMsgWithdrawValidatorRewardsAll struct {
	ValidatorAddr string `bson:"validator_addr"`
}

func (doctx *DocTxMsgWithdrawValidatorRewardsAll) Type() string {
	return constant.TxTypeWithdrawDelegatorRewardsAll
}

func (doctx *DocTxMsgWithdrawValidatorRewardsAll) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgWithdrawValidatorRewardsAll)
	doctx.ValidatorAddr = msg.ValidatorAddr.String()
}
