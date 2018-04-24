package collection

import (
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	DocsNmStakeTx = "tx_stake"
)

//Stake交易
type StakeTx struct {
	TxHash string    `bson:"tx_hash"`
	Time   time.Time `bson:"time"`
	Height int64     `bson:"height"`
	From   string    `bson:"from"`
	PubKey string    `bson:"pub_key"`
	Type   string    `bson:"type"`
	Amount coin.Coin `bson:"amount"`
}

func (c StakeTx) Name() string {
	return DocsNmStakeTx
}

func (c StakeTx) PkKvPair() map[string]interface{} {
	return bson.M{"tx_hash": c.TxHash}
}

func (c StakeTx) Index() mgo.Index {
	return mgo.Index{
		Key:        []string{"from"}, // 索引字段， 默认升序,若需降序在字段前加-
		Unique:     false,            // 唯一索引 同mysql唯一索引
		DropDups:   false,            // 索引重复替换旧文档,Unique为true时失效
		Background: true,             // 后台创建索引
	}
}
