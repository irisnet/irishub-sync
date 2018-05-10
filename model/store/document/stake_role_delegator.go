package document

import (
	"errors"
	"github.com/irisnet/iris-sync-server/model/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

const (
	CollectionNmStakeRoleDelegator = "stake_role_delegator"
)

type Delegator struct {
	Address string `bson:"address"`
	PubKey  string `bson:"pub_key"`
	Shares  int64  `bson:"shares"`
	UpdateTime time.Time `bson:"update_time"`
}

func (d Delegator) Name() string {
	return CollectionNmStakeRoleDelegator
}

func (d Delegator) PkKvPair() map[string]interface{} {
	return bson.M{"address": d.Address, "pub_key": d.PubKey}
}

func (d Delegator) Index() []mgo.Index {
	return []mgo.Index{
		{
			Key:        []string{"address"},
			Unique:     false,
			DropDups:   false,
			Background: true,
		},
		{
			Key:        []string{"pub_key"},
			Unique:     false,
			DropDups:   false,
			Background: true,
		},
		{
			Key:        []string{"address", "pub_key"},
			Unique:     true,
			DropDups:   false,
			Background: true,
		},
	}
}

func QueryDelegatorByAddressAndPubkey(address string, pubKey string) (Delegator, error) {
	var result Delegator
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{"address": address, "pub_key": pubKey}).Sort("-shares").One(&result)
		return err
	}

	if store.ExecCollection(CollectionNmStakeRoleDelegator, query) != nil {
		log.Printf("delegator is Empty")
		return result, errors.New("delegator is Empty")
	}

	return result, nil
}
