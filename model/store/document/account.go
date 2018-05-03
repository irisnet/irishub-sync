package document

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"github.com/irisnet/iris-sync-server/model/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

const (
	CollectionNmAccount = "account"
)

// 账户信息
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

//Account
func QueryAccount(address string) (Account, error) {
	var result Account
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{"address": address}).Sort("-amount").One(&result)
		return err
	}

	if store.ExecCollection(CollectionNmAccount, query) != nil {
		log.Printf("Account is Empry")
		return result, errors.New("Account is Empry")
	}

	return result, nil
}

func QueryAll() []Account {
	result := []Account{}
	query := func(c *mgo.Collection) error {
		err := c.Find(nil).All(&result)
		return err
	}

	if store.ExecCollection(CollectionNmAccount, query) != nil {
		log.Printf("Account is Empry")
	}
	return result
}
