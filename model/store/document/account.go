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

//账户信息
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

func (a Account) Index() mgo.Index {
	return mgo.Index{
		Key:        []string{"address"}, // 索引字段， 默认升序,若需降序在字段前加-
		Unique:     true,                // 唯一索引 同mysql唯一索引
		DropDups:   false,               // 索引重复替换旧文档,Unique为true时失效
		Background: true,                // 后台创建索引
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
