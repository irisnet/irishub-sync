package types

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/constant"
)

type SubmitProposal struct {
	Title          string      `json:"title"`          //  Title of the proposal
	Description    string      `json:"description"`    //  Description of the proposal
	Proposer       string      `json:"proposer"`       //  Address of the proposer
	InitialDeposit store.Coins `json:"initialDeposit"` //  Initial deposit paid by sender. Must be strictly positive.
	ProposalType   string      `json:"proposalType"`   //  Initial deposit paid by sender. Must be strictly positive.
	Params         Params      `json:"params"`
}

type SubmitSoftwareUpgradeProposal struct {
	SubmitProposal
	Version      uint64 `json:"version"`
	Software     string `json:"software"`
	SwitchHeight uint64 `json:"switch_height"`
	Threshold    string `json:"threshold"`
}

type SubmitTaxUsageProposal struct {
	SubmitProposal
	Usage       string `json:"usage"`
	DestAddress string `json:"dest_address"`
	Percent     string `json:"percent"`
}

type SubmitTokenAdditionProposal struct {
	SubmitProposal
	Symbol          string `json:"symbol"`
	CanonicalSymbol string `json:"canonical_symbol"`
	Name            string `json:"name"`
	Decimal         uint8  `json:"decimal"`
	MinUnitAlias    string `json:"min_unit_alias"`
	InitialSupply   uint64 `json:"initial_supply"`
}

type SetMemoRegexp struct {
	Owner      string `json:"owner"`
	MemoRegexp string `json:"memo_regexp"`
}

type RequestRand struct {
	Consumer      string `json:"consumer"`       // request address
	BlockInterval uint64 `json:"block-interval"` // block interval after which the requested random number will be generated
}

type Param struct {
	Subspace string `json:"subspace" bson:"subspace"`
	Key      string `json:"key" bson:"key"`
	Value    string `json:"value" bson:"value"`
}

type Params []Param

func NewSubmitProposal(msg MsgSubmitProposal) SubmitProposal {
	var params Params
	for _, p := range msg.Params {
		params = append(params, Param{
			Subspace: p.Subspace,
			Key:      p.Key,
			Value:    p.Value,
		})
	}
	return SubmitProposal{
		Title:          msg.Title,
		Description:    msg.Description,
		ProposalType:   msg.ProposalType.String(),
		Proposer:       msg.Proposer.String(),
		InitialDeposit: ParseCoins(msg.InitialDeposit.String()),
		Params:         params,
	}
}
func NewSubmitSoftwareUpgradeProposal(msg MsgSubmitSoftwareUpgradeProposal) SubmitSoftwareUpgradeProposal {
	submitProposal := NewSubmitProposal(msg.MsgSubmitProposal)
	return SubmitSoftwareUpgradeProposal{
		SubmitProposal: submitProposal,
		Version:        msg.Version,
		Software:       msg.Software,
		SwitchHeight:   msg.SwitchHeight,
		Threshold:      msg.Threshold.String(),
	}
}

func NewSubmitTaxUsageProposal(msg MsgSubmitTaxUsageProposal) SubmitTaxUsageProposal {
	submitProposal := NewSubmitProposal(msg.MsgSubmitProposal)
	return SubmitTaxUsageProposal{
		SubmitProposal: submitProposal,
		Usage:          msg.Usage.String(),
		DestAddress:    msg.DestAddress.String(),
		Percent:        msg.Percent.String(),
	}
}

func NewSubmitTokenAdditionProposal(msg MsgSubmitTokenAdditionProposal) SubmitTokenAdditionProposal {
	submitProposal := NewSubmitProposal(msg.MsgSubmitProposal)
	return SubmitTokenAdditionProposal{
		SubmitProposal: submitProposal,
		Symbol:msg.Symbol,
		CanonicalSymbol:msg.CanonicalSymbol,
		Name:msg.Name,
		Decimal:msg.Decimal,
		MinUnitAlias:msg.MinUnitAlias,
		InitialSupply:msg.InitialSupply,
	}
}
func NewRequestRand(msg MsgRequestRand)  RequestRand {
	return RequestRand{
		Consumer:msg.Consumer.String(),
		BlockInterval:msg.BlockInterval,
	}
}

func NewSetMemoRegexp(msg MsgSetMemoRegexp)  SetMemoRegexp {
	return SetMemoRegexp{
		Owner:msg.Owner.String(),
		MemoRegexp:msg.MemoRegexp,
	}
}

func (s SubmitProposal) Type() string {
	return constant.TxMsgTypeSubmitProposal
}

func (s SubmitProposal) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func (s SubmitSoftwareUpgradeProposal) Type() string {
	return constant.TxMsgTypeSubmitSoftwareUpgradeProposal
}

func (s SubmitSoftwareUpgradeProposal) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func (s SubmitTaxUsageProposal) Type() string {
	return constant.TxMsgTypeSubmitTaxUsageProposal
}

func (s SubmitTaxUsageProposal) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func (s SubmitTokenAdditionProposal) Type() string {
	return constant.TxMsgTypeSubmitTokenAdditionProposal
}

func (s SubmitTokenAdditionProposal) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func (s SetMemoRegexp) Type() string {
	return constant.TxTypeSetMemoRegexp
}
func (s SetMemoRegexp) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func (s RequestRand) Type() string {
	return constant.TxTypeRequestRand
}
func (s RequestRand) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func UnmarshalSubmitProposal(str string) (submitProposal SubmitProposal) {
	json.Unmarshal([]byte(str), &submitProposal)
	return
}

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

type Deposit struct {
	ProposalID uint64      `json:"proposal_id"` // ID of the proposal
	Depositer  string      `json:"depositer"`   // Address of the depositer
	Amount     store.Coins `json:"amount"`      // Coins to add to the proposal's deposit
}

func NewDeposit(deposit MsgDeposit) Deposit {
	return Deposit{
		ProposalID: deposit.ProposalID,
		Depositer:  deposit.Depositor.String(),
		Amount:     ParseCoins(deposit.Amount.String()),
	}
}

func (s Deposit) Type() string {
	return constant.TxTypeDeposit
}

func (s Deposit) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func UnmarshalDeposit(str string) (deposit Deposit) {
	json.Unmarshal([]byte(str), &deposit)
	return
}
