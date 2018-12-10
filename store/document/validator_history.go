package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const CollectionNmValidatorHistory = "validator_history"

type ValidatorHistory struct {
	Candidate
	UpdateTime time.Time `bson:"update_time"`
}

func (v ValidatorHistory) Name() string {
	return CollectionNmValidatorHistory
}

func (v ValidatorHistory) PkKvPair() map[string]interface{} {
	return bson.M{ValidatorUpTime_Field_ValAddress: v.Address}
}

func (v ValidatorHistory) RemoveAll() error {
	remove := func(c *mgo.Collection) error {
		_, err := c.RemoveAll(nil)
		return err
	}
	return store.ExecCollection(v.Name(), remove)
}

func (v ValidatorHistory) SaveAll(history []ValidatorHistory) error {
	var docs []interface{}

	if len(history) == 0 {
		return nil
	}

	for _, v := range history {
		docs = append(docs, v)
	}

	err := store.SaveAll(v.Name(), docs)

	return err
}

func (v ValidatorHistory) QueryAll() (vs []ValidatorHistory) {
	queryOp := func(c *mgo.Collection) error {
		return c.Find(nil).All(&vs)
	}
	store.ExecCollection(v.Name(), queryOp)
	return vs
}
