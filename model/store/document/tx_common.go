package document

import (
	"time"
	"github.com/irisnet/irishub-sync/model/store"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNmCommonTx = "tx_common"
)

type CommonTx struct {
	TxHash string     `bson:"tx_hash"`
	Time   time.Time  `bson:"time"`
	Height int64      `bson:"height"`
	From   string     `bson:"from"`
	To     string     `bson:"to"`
	Amount store.Coins `bson:"amount"`
	Type   string     `bson:"type"`
	Fee    store.Fee  `bson:"fee"`
	Status string `bson:"status"`
}

func (d CommonTx) Name() string {
	return CollectionNmCommonTx
}

func (d CommonTx) PkKvPair() map[string]interface{} {
	return bson.M{"tx_hash": d.TxHash}
}
