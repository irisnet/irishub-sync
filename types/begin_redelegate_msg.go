package types

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/util/constant"
)

type BeginRedelegate struct {
	DelegatorAddr    string  `json:"delegator_addr"`
	ValidatorSrcAddr string  `json:"validator_src_addr"`
	ValidatorDstAddr string  `json:"validator_dst_addr"`
	SharesAmount     float64 `json:"shares_amount"`
}

func (s BeginRedelegate) Type() string {
	return constant.TxTypeBeginRedelegate
}

func (s BeginRedelegate) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func NewBeginRedelegate(msg MsgBeginRedelegate) BeginRedelegate {
	shares, _ := msg.SharesAmount.Float64()
	return BeginRedelegate{
		DelegatorAddr:    msg.DelegatorAddr.String(),
		ValidatorSrcAddr: msg.ValidatorSrcAddr.String(),
		ValidatorDstAddr: msg.ValidatorDstAddr.String(),
		SharesAmount:     shares,
	}
}
