package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmCommonTx = "tx_common"
)

type CommonTx struct {
	Time   time.Time   `bson:"time"`
	Height int64       `bson:"height"`
	TxHash string      `bson:"tx_hash"`
	From   string      `bson:"from"`
	To     string      `bson:"to"`
	Amount store.Coins `bson:"amount"`
	Type   string      `bson:"type"`
	Fee    store.Fee   `bson:"fee"`
	Memo   string      `bson:"memo"`
	Status string      `bson:"status"`
	Log    string      `bson:"log"`

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
