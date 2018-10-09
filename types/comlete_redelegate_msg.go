package types

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/util/constant"
)

type CompleteRedelegate struct {
	DelegatorAddr    string `json:"delegator_addr"`
	ValidatorSrcAddr string `json:"validator_src_addr"`
	ValidatorDstAddr string `json:"validator_dst_addr"`
}

func (s CompleteRedelegate) Type() string {
	return constant.TxTypeBeginRedelegate
}

func (s CompleteRedelegate) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func NewCompleteRedelegate(msg MsgCompleteRedelegate) CompleteRedelegate {
	return CompleteRedelegate{
		DelegatorAddr:    msg.DelegatorAddr.String(),
		ValidatorSrcAddr: msg.ValidatorSrcAddr.String(),
		ValidatorDstAddr: msg.ValidatorDstAddr.String(),
	}
}
