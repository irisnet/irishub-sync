package types

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/util/constant"
)

type Vote struct {
	ProposalID uint64 `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     string `json:"option"`
}

func NewVote(v MsgVote) Vote {
	return Vote{
		ProposalID: v.ProposalID,
		Voter:      v.Voter.String(),
		Option:     v.Option.String(),
	}
}

func (s Vote) Type() string {
	return constant.TxTypeVote
}

func (s Vote) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func UnmarshalVote(str string) (vote Vote) {
	json.Unmarshal([]byte(str), &vote)
	return
}
