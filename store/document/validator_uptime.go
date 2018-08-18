package document

import (
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const CollectionName = "validator_uptime"

type ValidatorUpTime struct {
	ValAddress string  `bson:"val_address"`
	UpTime     float64 `bson:"up_time"`
}

func (d ValidatorUpTime) Name() string {
	return CollectionName
}

func (d ValidatorUpTime) PkKvPair() map[string]interface{} {
	return bson.M{"val_address": d.ValAddress}
}

func (d ValidatorUpTime) RemoveAll() error {
	query := bson.M{}
	remove := func(c *mgo.Collection) error {
		changeInfo, err := c.RemoveAll(query)
		logger.Info.Printf("remove all validator uptime data, %+v", changeInfo)
		return err
	}
	return store.ExecCollection(d.Name(), remove)
}

func (d ValidatorUpTime) SaveAll(validatorUpTimes []ValidatorUpTime) error {
	var docs []interface{}

	for _, v := range validatorUpTimes {
		docs = append(docs, v)
	}

	err := store.SaveAll(d.Name(), docs)

	return err
}
