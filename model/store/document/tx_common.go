package document

import (
	"time"
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
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
	Amount coin.Coins `bson:"amount"`
	Type   string     `bson:"type"`
}

func (d CommonTx) Name() string {
	return CollectionNmCommonTx
}

func (d CommonTx) PkKvPair() map[string]interface{} {
	return bson.M{"tx_hash": d.TxHash}
}

func (d CommonTx) Index() []mgo.Index {
	return []mgo.Index{}
}
