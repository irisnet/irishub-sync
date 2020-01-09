package msg

import (
	"github.com/irisnet/irishub-sync/store"
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub/app/v1/gov"
)

type DocTxMsgSubmitProposal struct {
	Title          string      `bson:"title"`          //  Title of the proposal
	Description    string      `bson:"description"`    //  Description of the proposal
	Proposer       string      `bson:"proposer"`       //  Address of the proposer
	InitialDeposit store.Coins `bson:"initialDeposit"` //  Initial deposit paid by sender. Must be strictly positive.
	ProposalType   string      `bson:"proposalType"`   //  Initial deposit paid by sender. Must be strictly positive.
	Params         Params      `bson:"params"`
}

type Param struct {
	Subspace string `json:"subspace" bson:"subspace"`
	Key      string `json:"key" bson:"key"`
	Value    string `json:"value" bson:"value"`
}

type Params []Param

func (doctx *DocTxMsgSubmitProposal) Type() string {
	return constant.TxTypeSubmitProposal
}

func (doctx *DocTxMsgSubmitProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgSubmitProposal)
	doctx.Title = msg.Title
	doctx.Description = msg.Description
	doctx.ProposalType = msg.ProposalType.String()
	doctx.Proposer = msg.Proposer.String()
	doctx.Params = loadParams(msg.Params)
	doctx.InitialDeposit = itypes.ParseCoins(msg.InitialDeposit.String())
}

func loadParams(params []gov.Param) (result []Param) {
	for _, val := range params {
		result = append(result, Param{Subspace: val.Subspace, Value: val.Value, Key: val.Key})
	}
	return
}

type DocTxMsgSubmitSoftwareUpgradeProposal struct {
	DocTxMsgSubmitProposal
	Version      uint64 `bson:"version"`
	Software     string `bson:"software"`
	SwitchHeight uint64 `bson:"switch_height"`
	Threshold    string `bson:"threshold"`
}

func (doctx *DocTxMsgSubmitSoftwareUpgradeProposal) Type() string {
	return constant.TxMsgTypeSubmitSoftwareUpgradeProposal
}

func (doctx *DocTxMsgSubmitSoftwareUpgradeProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgSubmitSoftwareUpgradeProposal)
	doctx.Title = msg.Title
	doctx.Description = msg.Description
	doctx.ProposalType = msg.ProposalType.String()
	doctx.Proposer = msg.Proposer.String()
	doctx.Params = loadParams(msg.Params)
	doctx.InitialDeposit = itypes.ParseCoins(msg.InitialDeposit.String())
	doctx.Version = msg.Version
	doctx.Software = msg.Software
	doctx.SwitchHeight = msg.SwitchHeight
	doctx.Threshold = msg.Threshold.String()
}

type DocTxMsgSubmitCommunityTaxUsageProposal struct {
	DocTxMsgSubmitProposal
	Usage       string `bson:"usage"`
	DestAddress string `bson:"dest_address"`
	Percent     string `bson:"percent"`
}

func (doctx *DocTxMsgSubmitCommunityTaxUsageProposal) Type() string {
	return constant.TxMsgTypeSubmitTaxUsageProposal
}

func (doctx *DocTxMsgSubmitCommunityTaxUsageProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgSubmitTaxUsageProposal)
	doctx.Title = msg.Title
	doctx.Description = msg.Description
	doctx.ProposalType = msg.ProposalType.String()
	doctx.Proposer = msg.Proposer.String()
	doctx.Params = loadParams(msg.Params)
	doctx.InitialDeposit = itypes.ParseCoins(msg.InitialDeposit.String())
	doctx.Usage = msg.Usage.String()
	doctx.DestAddress = msg.DestAddress.String()
	doctx.Percent = msg.Percent.String()
}

type DocTxMsgSubmitTokenAdditionProposal struct {
	DocTxMsgSubmitProposal
	Symbol          string `bson:"symbol"`
	CanonicalSymbol string `bson:"canonical_symbol"`
	Name            string `bson:"name"`
	Decimal         uint8  `bson:"decimal"`
	MinUnitAlias    string `bson:"min_unit_alias"`
	//InitialSupply   uint64 `bson:"initial_supply"`
}

func (doctx *DocTxMsgSubmitTokenAdditionProposal) Type() string {
	return constant.TxMsgTypeSubmitTokenAdditionProposal
}

func (doctx *DocTxMsgSubmitTokenAdditionProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgSubmitTokenAdditionProposal)
	var params Params
	for _, p := range msg.Params {
		params = append(params, Param{
			Subspace: p.Subspace,
			Key:      p.Key,
			Value:    p.Value,
		})
	}
	doctx.Title = msg.Title
	doctx.Description = msg.Description
	doctx.ProposalType = msg.ProposalType.String()
	doctx.Proposer = msg.Proposer.String()
	doctx.Params = params
	doctx.InitialDeposit = itypes.ParseCoins(msg.InitialDeposit.String())
	doctx.Symbol = msg.Symbol
	doctx.MinUnitAlias = msg.MinUnitAlias
	doctx.CanonicalSymbol = msg.CanonicalSymbol
	doctx.Name = msg.Name
	doctx.Decimal = msg.Decimal
}

// MsgVote
type DocTxMsgVote struct {
	ProposalID uint64 `bson:"proposal_id"` // ID of the proposal
	Voter      string `bson:"voter"`       //  address of the voter
	Option     string `bson:"option"`      //  option from OptionSet chosen by the voter
}

func (doctx *DocTxMsgVote) Type() string {
	return constant.TxTypeVote
}

func (doctx *DocTxMsgVote) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgVote)
	doctx.Voter = msg.Voter.String()
	doctx.Option = msg.Option.String()
	doctx.ProposalID = msg.ProposalID
}

// MsgDeposit
type DocTxMsgDeposit struct {
	ProposalID uint64      `bson:"proposal_id"` // ID of the proposal
	Depositor  string      `bson:"depositor"`   // Address of the depositor
	Amount     store.Coins `bson:"amount"`      // Coins to add to the proposal's deposit
}

func (doctx *DocTxMsgDeposit) Type() string {
	return constant.TxTypeDeposit
}

func (doctx *DocTxMsgDeposit) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgDeposit)
	doctx.Depositor = msg.Depositor.String()
	doctx.Amount = itypes.ParseCoins(msg.Amount.String())
	doctx.ProposalID = msg.ProposalID
}
