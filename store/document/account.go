package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNmAccount = "account"
	AccountFieldAddress = "address"
)

type Account struct {
	Address          string     `bson:"address"`
	AccountNumber    uint64     `bson:"account_number"`
	CoinIris         store.Coin `bson:"coin_iris"`
	CoinIrisUpdateAt int64      `bson:"coin_iris_update_at"`
}

func (a Account) Name() string {
	return CollectionNmAccount
}

func (a Account) PkKvPair() map[string]interface{} {
	return bson.M{AccountFieldAddress: a.Address}
}
