package document

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/constant"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmCommonTx = "tx_common"
)

type CommonTx struct {
	Time      time.Time       `bson:"time"`
	Height    int64           `bson:"height"`
	TxHash    string          `bson:"tx_hash"`
	From      string          `bson:"from"`
	To        string          `bson:"to"`
	Amount    store.Coins     `bson:"amount"`
	Type      string          `bson:"type"`
	Fee       store.Fee       `bson:"fee"`
	Memo      string          `bson:"memo"`
	Status    string          `bson:"status"`
	Log       string          `bson:"log"`
	GasUsed   int64           `bson:"gas_used"`
	GasPrice  float64         `bson:"gas_price"`
	ActualFee store.ActualFee `bson:"actual_fee"`

	StakeCreateValidator StakeCreateValidator `bson:"stake_create_validator"`
	StakeEditValidator   StakeEditValidator   `bson:"stake_edit_validator"`
}

// Description
type ValDescription struct {
	Moniker  string `bson:"moniker"`
	Identity string `bson:"identity"`
	Website  string `bson:"website"`
	Details  string `bson:"details"`
}

type StakeCreateValidator struct {
	PubKey      string         `bson:"pub_key"`
	Description ValDescription `bson:"description"`
}

type StakeEditValidator struct {
	Description ValDescription `bson:"description"`
}

func (d CommonTx) Name() string {
	return CollectionNmCommonTx
}

func (d CommonTx) PkKvPair() map[string]interface{} {
	return bson.M{"tx_hash": d.TxHash}
}

func (d CommonTx) Query(query, fields bson.M, sort []string, skip, limit int) (
	results []CommonTx, err error) {
	exop := func(c *mgo.Collection) error {
		return c.Find(query).Sort(sort...).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	return results, store.ExecCollection(d.Name(), exop)
}

func (d CommonTx) CalculateTxGasAndGasPrice(txType string, limit int) (
	[]CommonTx, error) {
	query := bson.M{
		"type":   txType,
		"status": constant.TxStatusSuccess,
	}
	fields := bson.M{}
	sort := []string{"-height"}
	skip := 0

	return d.Query(query, fields, sort, skip, limit)
}
