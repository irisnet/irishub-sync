package document

import (
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmCoinTx = "tx_coin"
)

//Coin交易
type CoinTx struct {
	TxHash string     `bson:"tx_hash"`
	Time   time.Time  `bson:"time"`
	Height int64      `bson:"height"`
	From   string     `bson:"from"`
	To     string     `bson:"to"`
	Amount coin.Coins `bson:"amount"`
}

func (c CoinTx) Name() string {
	return CollectionNmCoinTx
}

func (c CoinTx) PkKvPair() map[string]interface{} {
	return bson.M{"tx_hash": c.TxHash}
}

func (c CoinTx) Index() mgo.Index {
	return mgo.Index{
		Key:        []string{"from", "to"}, // 索引字段， 默认升序,若需降序在字段前加-
		Unique:     false,                  // 唯一索引 同mysql唯一索引
		DropDups:   false,                  // 索引重复替换旧文档,Unique为true时失效
		Background: true,                   // 后台创建索引
	}
}
