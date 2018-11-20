package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const CollectionNmProposal = "proposal"

type Proposal struct {
	ProposalId      uint64      `bson:"proposal_id"`
	Title           string      `bson:"title"`
	Type            string      `bson:"type"`
	Description     string      `bson:"description"`
	Status          string      `bson:"status"`
	SubmitTime      time.Time   `bson:"submit_time"`
	DepositEndTime  time.Time   `bson:"deposit_end_time"`
	VotingStartTime time.Time   `bson:"voting_start_time"`
	VotingEndTime   time.Time   `bson:"voting_end_time"`
	TotalDeposit    store.Coins `bson:"total_deposit"`
	Votes           []PVote     `bson:"votes"`
}

type PVote struct {
	Voter  string    `json:"voter"`
	Option string    `json:"option"`
	Time   time.Time `json:"time"`
}

func (m Proposal) Name() string {
	return CollectionNmProposal
}

func (m Proposal) PkKvPair() map[string]interface{} {
	return bson.M{"proposal_id": m.ProposalId}
}

func QueryProposal(proposalId uint64) (Proposal, error) {
	var result Proposal
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{"proposal_id": proposalId}).Sort("-submit_block").One(&result)
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
		err := c.Find(bson.M{"status": bson.M{"$in": status}}).All(&result)
		return err
	}
	err := store.ExecCollection(CollectionNmProposal, query)

	if err != nil {
		return result, err
	}
	return result, nil
}
