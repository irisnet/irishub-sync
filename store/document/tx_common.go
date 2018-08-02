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
}

func (d CommonTx) Name() string {
	return CollectionNmCommonTx
}

func (d CommonTx) PkKvPair() map[string]interface{} {
	return bson.M{"tx_hash": d.TxHash}
}
