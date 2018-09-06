package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const CollectionNmProposal = "proposal"

type Proposal struct {
	ProposalId       int64       `bson:"proposal_id"`
	Title            string      `bson:"title"`
	Type             string      `bson:"type"`
	Description      string      `bson:"description"`
	Status           string      `bson:"status"`
	SubmitBlock      int64       `bson:"submit_block"`
	SubmitTime       time.Time   `bson:"submit_time"`
	VotingStartBlock int64       `bson:"voting_start_block"`
	TotalDeposit     store.Coins `bson:"total_deposit"`
	Votes            []PVote     `bson:"votes"`
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

func QueryProposal(proposalId int64) (Proposal, error) {
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
