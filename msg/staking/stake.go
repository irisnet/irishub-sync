package staking

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	stake "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// MsgDelegate - struct for bonding transactions
type DocTxMsgBeginRedelegate struct {
	DelegatorAddr    string `bson:"delegator_addr"`
	ValidatorSrcAddr string `bson:"validator_src_addr"`
	ValidatorDstAddr string `bson:"validator_dst_addr"`
	SharesAmount     string `bson:"shares_amount"`
}

func (doctx *DocTxMsgBeginRedelegate) Type() string {
	return constant.TxTypeBeginRedelegate
}

func (doctx *DocTxMsgBeginRedelegate) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgBeginRedelegate)
	doctx.DelegatorAddr = msg.DelegatorAddress.String()
	doctx.ValidatorSrcAddr = msg.ValidatorSrcAddress.String()
	doctx.ValidatorDstAddr = msg.ValidatorDstAddress.String()
	doctx.SharesAmount = msg.Amount.String()
}

// MsgUnjail - struct for unjailing jailed validator
type DocTxMsgUnjail struct {
	ValidatorAddr string `bson:"address"` // address of the validator operator
}

func (doctx *DocTxMsgUnjail) Type() string {
	return constant.TxTypeUnjail
}

func (doctx *DocTxMsgUnjail) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgUnjail)
	doctx.ValidatorAddr = msg.ValidatorAddr.String()
}

// MsgBeginUnbonding - struct for unbonding transactions
type DocTxMsgBeginUnbonding struct {
	DelegatorAddr string `bson:"delegator_addr"`
	ValidatorAddr string `bson:"validator_addr"`
	SharesAmount  string `bson:"shares_amount"`
}

func (doctx *DocTxMsgBeginUnbonding) Type() string {
	return constant.TxTypeStakeBeginUnbonding
}

func (doctx *DocTxMsgBeginUnbonding) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgStakeBeginUnbonding)
	doctx.ValidatorAddr = msg.ValidatorAddress.String()
	doctx.DelegatorAddr = msg.DelegatorAddress.String()
	doctx.SharesAmount = msg.Amount.String()
}

// MsgDelegate - struct for bonding transactions
type DocTxMsgDelegate struct {
	DelegatorAddr string     `bson:"delegator_addr"`
	ValidatorAddr string     `bson:"validator_addr"`
	Delegation    store.Coin `bson:"delegation"`
}

func (doctx *DocTxMsgDelegate) Type() string {
	return constant.TxTypeStakeDelegate
}

func (doctx *DocTxMsgDelegate) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgStakeDelegate)
	doctx.ValidatorAddr = msg.ValidatorAddress.String()
	doctx.DelegatorAddr = msg.DelegatorAddress.String()
	doctx.Delegation = itypes.ParseCoin(msg.Amount.String())
}

// MsgEditValidator - struct for editing a validator
type DocTxMsgStakeEdit struct {
	document.ValDescription
	ValidatorAddr  string `bson:"address"`
	CommissionRate string `bson:"commission_rate"`
}

func (doctx *DocTxMsgStakeEdit) Type() string {
	return constant.TxTypeStakeEditValidator
}

func (doctx *DocTxMsgStakeEdit) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgStakeEdit)
	doctx.ValidatorAddr = msg.ValidatorAddress.String()
	commissionRate := msg.CommissionRate
	if commissionRate == nil {
		doctx.CommissionRate = ""
	} else {
		doctx.CommissionRate = commissionRate.String()
	}
	doctx.ValDescription = loadDescription(msg.Description)
}

type DocTxMsgStakeCreate struct {
	document.ValDescription
	Commission    document.CommissionMsg
	DelegatorAddr string     `bson:"delegator_address"`
	ValidatorAddr string     `bson:"validator_address"`
	PubKey        string     `bson:"pubkey"`
	Delegation    store.Coin `bson:"delegation"`
}

func (doctx *DocTxMsgStakeCreate) Type() string {
	return constant.TxTypeStakeCreateValidator
}

func (doctx *DocTxMsgStakeCreate) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgStakeCreate)
	//pubKey, err := itypes.Bech32ifyValPub(msg.Pubkey)
	//if err != nil {
	//	pubKey = ""
	//}
	doctx.ValidatorAddr = msg.ValidatorAddress.String()
	doctx.PubKey = msg.Pubkey
	doctx.DelegatorAddr = msg.DelegatorAddress.String()
	doctx.Delegation = itypes.ParseCoin(msg.Description.String())
	doctx.Commission = document.CommissionMsg{
		Rate:          msg.Commission.Rate.String(),
		MaxChangeRate: msg.Commission.MaxChangeRate.String(),
		MaxRate:       msg.Commission.MaxRate.String(),
	}
	doctx.ValDescription = loadDescription(msg.Description)
}

func loadDescription(description stake.Description) document.ValDescription {
	return document.ValDescription{
		Moniker:  description.Moniker,
		Details:  description.Details,
		Identity: description.Identity,
		Website:  description.Website,
	}
}
