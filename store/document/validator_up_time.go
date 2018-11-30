package document

import (
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionName = "validator_up_time"

	ValidatorUpTime_Field_ValAddress = "val_address"
	ValidatorUpTime_Field_UpTime     = "up_time"
)

type ValidatorUpTime struct {
	ValAddress string  `bson:"val_address"`
	UpTime     float64 `bson:"up_time"`
}

func (d ValidatorUpTime) Name() string {
	return CollectionName
}

func (d ValidatorUpTime) PkKvPair() map[string]interface{} {
	return bson.M{ValidatorUpTime_Field_ValAddress: d.ValAddress}
}

func (d ValidatorUpTime) RemoveAll() error {
	query := bson.M{}
	remove := func(c *mgo.Collection) error {
		changeInfo, err := c.RemoveAll(query)
		logger.Info("remove all validator uptime data", logger.Any("changeInfo", changeInfo))
		return err
	}
	return store.ExecCollection(d.Name(), remove)
}

func (d ValidatorUpTime) SaveAll(validatorUpTimes []ValidatorUpTime) error {
	var docs []interface{}

	if len(validatorUpTimes) == 0 {
		return nil
	}

	for _, v := range validatorUpTimes {
		docs = append(docs, v)
	}

	err := store.SaveAll(d.Name(), docs)

	return err
}
