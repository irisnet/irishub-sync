package document

import (
	"errors"
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmStakeRoleDelegator = "stake_role_delegator"
)

type Delegator struct {
	Address       string    `bson:"address"`
	ValidatorAddr string    `bson:"validator_addr"` // validatorAddr
	Shares        int64     `bson:"shares"`
	UpdateTime    time.Time `bson:"update_time"`
}

func (d Delegator) Name() string {
	return CollectionNmStakeRoleDelegator
}

func (d Delegator) PkKvPair() map[string]interface{} {
	return bson.M{"address": d.Address, "validator_addr": d.ValidatorAddr}
}

func QueryDelegatorByAddressAndValAddr(address string, valAddr string) (Delegator, error) {
	var result Delegator
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{"address": address, "validator_addr": valAddr}).Sort("-shares").One(&result)
		return err
	}

	if store.ExecCollection(CollectionNmStakeRoleDelegator, query) != nil {
		return result, errors.New("delegator is Empty")
	}

	return result, nil
}
