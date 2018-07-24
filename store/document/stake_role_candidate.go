package document

import (
	"errors"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNmStakeRoleCandidate = "stake_role_candidate"
)

type Candidate struct {
	Address     string      `bson:"address"` // owner, identity key
	PubKey      string      `bson:"pub_key"`
	Revoked     bool        `bson:"revoked"` // has the validator been revoked from bonded status
	Shares      float64     `bson:"shares"`
	OriginalShares string   `bson:"original_shares"`
	VotingPower float64     `bson:"voting_power"` // Voting power if pubKey is a considered a validator
	Description Description `bson:"description"`  // Description terms for the candidate
	BondHeight  int64       `bson:"bond_height"`
}

func (d Candidate) Name() string {
	return CollectionNmStakeRoleCandidate
}

func (d Candidate) PkKvPair() map[string]interface{} {
	return bson.M{"address": d.Address}
}

func QueryCandidateByAddress(address string) (Candidate, error) {
	var result Candidate
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{"address": address}).One(&result)
		return err
	}

	if store.ExecCollection(CollectionNmStakeRoleCandidate, query) != nil {
		logger.Info.Println("candidate is Empty")
		return result, errors.New("candidate is Empty")
	}

	return result, nil
}
