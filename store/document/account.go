package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmAccount = "account"
)

type Account struct {
	Address string      `bson:"address"`
	Amount  store.Coins `bson:"amount"`
	Time    time.Time   `bson:"time"`
	Height  int64       `bson:"height"`
}

func (a Account) Name() string {
	return CollectionNmAccount
}

func (a Account) PkKvPair() map[string]interface{} {
	return bson.M{"address": a.Address}
}

func QueryAccount(address string) (Account, error) {
	var result Account
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{"address": address}).Sort("-amount.amount").One(&result)
		return err
	}

	err := store.ExecCollection(CollectionNmAccount, query)

	if err != nil {
		return result, err
	}

	return result, nil
}
