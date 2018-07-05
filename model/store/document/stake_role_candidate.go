package document

import (
	"errors"
	"github.com/irisnet/irishub-sync/model/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/irisnet/irishub-sync/module/logger"
)

const (
	CollectionNmStakeRoleCandidate = "stake_role_candidate"
)

type Candidate struct {
	Address     string      `bson:"address"` // owner, identity key
	PubKey      string      `bson:"pub_key"`
	Shares      int64       `bson:"shares"`
	VotingPower int64       `bson:"voting_power"` // Voting power if pubKey is a considered a validator
	Description Description `bson:"description"`  // Description terms for the candidate
	UpdateTime  time.Time   `bson:"update_time"`
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
