package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmProposal = "proposal"

	Proposal_Field_ProposalId = "proposal_id"
	Proposal_Field_Status     = "status"
)

type Proposal struct {
	ProposalId              uint64       `bson:"proposal_id"`
	Title                   string       `bson:"title"`
	Type                    string       `bson:"type"`
	Description             string       `bson:"description"`
	Status                  string       `bson:"status"`
	SubmitTime              time.Time    `bson:"submit_time"`
	DepositEndTime          time.Time    `bson:"deposit_end_time"`
	VotingStartTime         time.Time    `bson:"voting_start_time"`
	VotingEndTime           time.Time    `bson:"voting_end_time"`
	VotingPeriodStartHeight int64        `bson:"voting_start_height"`
	VotingEndHeight         int64        `bson:"voting_end_height"`
	TotalDeposit            store.Coins  `bson:"total_deposit"`
	Votes                   []PVote      `bson:"votes"`
	TallyResult             PTallyResult `bson:"tally_result"`
}

type PVote struct {
	Voter  string    `json:"voter" bson:"voter"`
	Option string    `json:"option" bson:"option"`
	TxHash string    `json:"tx_hash" bson:"tx_hash"`
	Time   time.Time `json:"time" bson:"time"`
}

//-----------------------------------------------------------
// Tally Results
type PTallyResult struct {
	Yes        string `json:"yes" bson:"yes"`
	Abstain    string `json:"abstain" bson:"abstain"`
	No         string `json:"no" bson:"no"`
	NoWithVeto string `json:"no_with_veto" bson:"no_with_veto"`
	SystemVotingPower string `json:"system_voting_power" bson:"system_voting_power"`
}

func (m Proposal) Name() string {
	return CollectionNmProposal
}

func (m Proposal) PkKvPair() map[string]interface{} {
	return bson.M{Proposal_Field_ProposalId: m.ProposalId}
}

func QueryProposal(proposalId uint64) (Proposal, error) {
	var result Proposal
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{Proposal_Field_ProposalId: proposalId}).Sort("-submit_block").One(&result)
		return err
	}

	err := store.ExecCollection(CollectionNmProposal, query)

	if err != nil {
		return result, err
	}

	return result, nil
}
func QueryByStatus(status []string) ([]Proposal, error) {
	var result []Proposal
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{Proposal_Field_Status: bson.M{"$in": status}}).All(&result)
		return err
	}
	err := store.ExecCollection(CollectionNmProposal, query)

	if err != nil {
		return result, err
	}
	return result, nil
}
