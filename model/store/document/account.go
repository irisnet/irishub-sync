package document

import (
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"github.com/irisnet/iris-sync-server/model/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/irisnet/iris-sync-server/module/logger"
)

const (
	CollectionNmAccount = "account"
)

type Account struct {
	Address string     `bson:"address"`
	Amount  coin.Coins `bson:"amount"`
	Time    time.Time  `bson:"time"`
	Height  int64      `bson:"height"`
}

func (a Account) Name() string {
	return CollectionNmAccount
}

func (a Account) PkKvPair() map[string]interface{} {
	return bson.M{"address": a.Address}
}

func (a Account) Index() []mgo.Index {
	return []mgo.Index{
		{
			Key:        []string{"address"},
			Unique:     true,
			DropDups:   false,
			Background: true,
		},
	}
}

func QueryAccount(address string) (Account, error) {
	var result Account
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{"address": address}).Sort("-amount.amount").One(&result)
		return err
	}

	if store.ExecCollection(CollectionNmAccount, query) != nil {
		logger.Info.Println("Account is Empty")
		return result, nil
	}

	return result, nil
}
