package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmAccount = "account"
	AccountFieldAddress = "address"
)

type Account struct {
	Address             string     `bson:"address"`
	AccountNumber       uint64     `bson:"account_number"`
	Total               store.Coin `bson:"total"`
	CoinIris            store.Coin `bson:"coin_iris"`
	Delegation          store.Coin `bson:"delegation"`
	UnbondingDelegation store.Coin `bson:"unbonding_delegation"`
	Rewards             store.Coin `bson:"rewards"`
	UpdateAt            int64      `bson:"update_at"`
	CreateAt            int64      `bson:"create_at"`
}

func (d Account) Name() string {
	return CollectionNmAccount
}

func (d Account) PkKvPair() map[string]interface{} {
	return bson.M{AccountFieldAddress: d.Address}
}

// override store.Save()
// not to check record if exist before save document
func (d Account) Save(account Account) error {
	account.CreateAt = time.Now().Unix()
	fn := func(c *mgo.Collection) error {
		return c.Insert(account)
	}

	return store.ExecCollection(d.Name(), fn)
}

// get account by primary key
// return a empty struct when record is not exists
func (d Account) getAccountByPK() (Account, error) {
	var (
		res Account
	)
	find := func(c *mgo.Collection) error {
		return c.Find(d.PkKvPair()).One(&res)
	}

	if err := store.ExecCollection(d.Name(), find); err != nil {
		if err == mgo.ErrNotFound {
			return res, nil
		} else {
			return res, err
		}
	}

	return res, nil
}

// save  account address
func (d Account) SaveAddress(address string) error {
	d.Address = address
	if account, err := d.getAccountByPK(); err != nil {
		return err
	} else {
		if account.Address != "" {
			return nil
		}
		account.Address = address
		return d.Save(account)

	}
}
