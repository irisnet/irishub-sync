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
	if err := store.Find(d.Name(), d.PkKvPair()).One(&res); err != nil {
		if err == mgo.ErrNotFound {
			return res, nil
		} else {
			return res, err
		}
	}

	return res, nil
}

// save or update account balance info
func (d Account) UpsertBalanceInfo(address string, balance store.Coin, accountNumber uint64, height, timestamp int64) error {
	d.Address = address
	if account, err := d.getAccountByPK(); err != nil {
		return err
	} else {
		account.Address = address
		account.AccountNumber = accountNumber
		account.CoinIris = balance
		account.CoinIrisUpdateHeight = height
		account.CoinIrisUpdateAt = timestamp
		account.TotalUpdateHeight = height
		account.TotalUpdateAt = timestamp

		if account.Address == "" {
			// record not exist
			account.Total = balance
			return d.Save(account)
		} else {
			// record already exist
			account.Total = store.Coin{
				Denom:  balance.Denom,
				Amount: balance.Amount + account.Delegation.Amount + account.UnbondingDelegation.Amount,
			}
			return store.Update(account)
		}
	}
}

// save or update delegation info
func (d Account) UpsertDelegationInfo(address string, delegation store.Coin, height, timestamp int64) error {
	d.Address = address
	if account, err := d.getAccountByPK(); err != nil {
		return err
	} else {
		account.Address = address
		account.Delegation = delegation
		account.DelegationUpdateHeight = height
		account.DelegationUpdateAt = timestamp
		account.TotalUpdateHeight = height
		account.TotalUpdateAt = timestamp
		if account.Address == "" {
			// record not exist
			account.Total = delegation
			return d.Save(account)
		} else {
			// record exist
			account.Total = store.Coin{
				Denom:  delegation.Denom,
				Amount: account.CoinIris.Amount + delegation.Amount + account.UnbondingDelegation.Amount,
			}
			return store.Update(account)
		}
	}
}

// save or update unbondingDelegation info
func (d Account) UpsertUnbondingDelegationInfo(address string, unbondingDelegation store.Coin, height, timestamp int64) error {
	d.Address = address
	if account, err := d.getAccountByPK(); err != nil {
		return err
	} else {
		account.Address = address
		account.UnbondingDelegation = unbondingDelegation
		account.UnbondingDelegationUpdateHeight = height
		account.UnbondingDelegationUpdateAt = timestamp
		account.TotalUpdateHeight = height
		account.TotalUpdateAt = timestamp
		if account.Address == "" {
			// record not exist
			account.Total = unbondingDelegation
			return d.Save(account)
		} else {
			// record exist
			account.Total = store.Coin{
				Denom:  unbondingDelegation.Denom,
				Amount: account.CoinIris.Amount + account.Delegation.Amount + unbondingDelegation.Amount,
			}
			return store.Update(account)
		}
	}
}
