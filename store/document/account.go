package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNmAccount = "account"
	AccountFieldAddress = "address"
)

type Account struct {
	Address                         string     `bson:"address"`
	AccountNumber                   uint64     `bson:"account_number"`
	CoinIris                        store.Coin `bson:"coin_iris"`
	Delegation                      store.Coin `bson:"delegation"`
	UnbondingDelegation             store.Coin `bson:"unbonding_delegation"`
	Total                           store.Coin `bson:"total"`
	CoinIrisUpdateHeight            int64      `bson:"coin_iris_update_height"`
	CoinIrisUpdateAt                int64      `bson:"coin_iris_update_at"`
	DelegationUpdateHeight          int64      `bson:"delegation_update_height"`
	DelegationUpdateAt              int64      `bson:"delegation_update_at"`
	UnbondingDelegationUpdateHeight int64      `bson:"unbonding_delegation_update_height"`
	UnbondingDelegationUpdateAt     int64      `bson:"unbonding_delegation_update_at"`
	TotalUpdateHeight               int64      `bson:"total_update_height"`
	TotalUpdateAt                   int64      `bson:"total_update_at"`
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

// save  account balance info
func (d Account) SaveBalanceInfo(address string, balance store.Coin, accountNumber uint64, height, timestamp int64) error {
	d.Address = address
	if account, err := d.getAccountByPK(); err != nil {
		return err
	} else {
		if account.Address != "" {
			return nil
		}
		// record not exist
		account.AccountNumber = accountNumber
		account.CoinIris = balance
		account.CoinIrisUpdateHeight = height
		account.CoinIrisUpdateAt = timestamp
		account.TotalUpdateHeight = height
		account.TotalUpdateAt = timestamp
		account.Address = address
		account.Total = balance
		return d.Save(account)

	}
}

// save  delegation info
func (d Account) SaveDelegationInfo(address string, delegation store.Coin, height, timestamp int64) error {
	d.Address = address
	if account, err := d.getAccountByPK(); err != nil {
		return err
	} else {
		if account.Address != "" {
			return nil
		}
		// record not exist
		account.Delegation = delegation
		account.DelegationUpdateHeight = height
		account.DelegationUpdateAt = timestamp
		account.TotalUpdateHeight = height
		account.TotalUpdateAt = timestamp
		account.Address = address
		account.Total = delegation
		return d.Save(account)

	}
}

// save unbondingDelegation info
func (d Account) SaveUnbondingDelegationInfo(address string, unbondingDelegation store.Coin, height, timestamp int64) error {
	d.Address = address
	if account, err := d.getAccountByPK(); err != nil {
		return err
	} else {

		if account.Address != "" {
			return nil
		}
		// record not exist
		account.UnbondingDelegation = unbondingDelegation
		account.UnbondingDelegationUpdateHeight = height
		account.UnbondingDelegationUpdateAt = timestamp
		account.TotalUpdateHeight = height
		account.TotalUpdateAt = timestamp
		account.Address = address
		account.Total = unbondingDelegation
		return d.Save(account)

	}
}
