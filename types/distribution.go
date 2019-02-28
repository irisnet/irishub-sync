package types

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/util/constant"
)

type SetWithdrawAddressMsg struct {
	DelegatorAddr string `json:"delegator_addr"`
	WithdrawAddr  string `json:"delegator_addr"`
}

func NewSetWithdrawAddressMsg(msg MsgSetWithdrawAddress) SetWithdrawAddressMsg {
	return SetWithdrawAddressMsg{
		DelegatorAddr: msg.DelegatorAddr.String(),
		WithdrawAddr:  msg.WithdrawAddr.String(),
	}
}

func (s SetWithdrawAddressMsg) Type() string {
	return constant.TxTypeSetWithdrawAddress
}

func (s SetWithdrawAddressMsg) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

type WithdrawDelegatorRewardsAllMsg struct {
	DelegatorAddr string `json:"delegator_addr"`
}

func NewWithdrawDelegatorRewardsAllMsg(msg MsgWithdrawDelegatorRewardsAll) WithdrawDelegatorRewardsAllMsg {
	return WithdrawDelegatorRewardsAllMsg{
		DelegatorAddr: msg.DelegatorAddr.String(),
	}
}

func (s WithdrawDelegatorRewardsAllMsg) Type() string {
	return constant.TxTypeWithdrawDelegatorRewardsAll
}

func (s WithdrawDelegatorRewardsAllMsg) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

type WithdrawDelegatorRewardMsg struct {
	DelegatorAddr string `json:"delegator_addr"`
	ValidatorAddr string `json:"validator_addr"`
}

func NewWithdrawDelegatorRewardMsg(msg MsgWithdrawDelegatorReward) WithdrawDelegatorRewardMsg {
	return WithdrawDelegatorRewardMsg{
		DelegatorAddr: msg.DelegatorAddr.String(),
		ValidatorAddr: msg.ValidatorAddr.String(),
	}
}

func (s WithdrawDelegatorRewardMsg) Type() string {
	return constant.TxTypeWithdrawDelegatorReward
}

func (s WithdrawDelegatorRewardMsg) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

type WithdrawValidatorRewardsAllMsg struct {
	ValidatorAddr string `json:"validator_addr"`
}

func NewWithdrawValidatorRewardsAllMsg(msg MsgWithdrawValidatorRewardsAll) WithdrawValidatorRewardsAllMsg {
	return WithdrawValidatorRewardsAllMsg{
		ValidatorAddr: msg.ValidatorAddr.String(),
	}
}

func (s WithdrawValidatorRewardsAllMsg) Type() string {
	return constant.TxTypeWithdrawValidatorRewardsAll
}

func (s WithdrawValidatorRewardsAllMsg) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}
