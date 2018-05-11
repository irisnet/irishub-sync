package document

import (
	"errors"
	"github.com/irisnet/iris-sync-server/model/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/irisnet/iris-sync-server/module/logger"
)

const (
	CollectionNmStakeRoleCandidate = "stake_role_candidate"
)

type Candidate struct {
	Address     string      `bson:"address"` // owner
	PubKey      string      `bson:"pub_key"`
	Shares      int64       `bson:"shares"`
	VotingPower uint64      `bson:"voting_power"` // Voting power if pubKey is a considered a validator
	Description Description `bson:"description"`  // Description terms for the candidate
	UpdateTime time.Time    `bson:"update_time"`
}

func (d Candidate) Name() string {
	return CollectionNmStakeRoleCandidate
}

func (d Candidate) PkKvPair() map[string]interface{} {
	return bson.M{"pub_key": d.PubKey}
}

func (d Candidate) Index() []mgo.Index {
	return []mgo.Index {
		{
			Key:        []string{"pub_key"},
			Unique:     true,
			DropDups:   false,
			Background: true,
		},
		{
			Key:        []string{"address"},
			Unique:     false,
			DropDups:   false,
			Background: true,
		},
	}
}

func QueryCandidateByPubkey(pubKey string) (Candidate, error) {
	var result Candidate
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{"pub_key": pubKey}).One(&result)
		return err
	}

	if store.ExecCollection(CollectionNmStakeRoleCandidate, query) != nil {
		logger.Info.Println("candidate is Empty")
		return result, errors.New("candidate is Empty")
	}

	return result, nil
}
