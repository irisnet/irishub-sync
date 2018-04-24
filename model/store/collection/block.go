package collection

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	DocsNmBlock = "block"
)

type Block struct {
	Height int64     `bson:"height"`
	Time   time.Time `bson:"time"`
	TxNum  int64     `bson:"tx_num"`
}

func (d Block) Name() string {
	return DocsNmBlock
}

func (d Block) PkKvPair() map[string]interface{} {
	return bson.M{"height": d.Height}
}

func (d Block) Index() mgo.Index {
	return mgo.Index{
		Key:        []string{"height"}, // 索引字段， 默认升序,若需降序在字段前加-
		Unique:     true,               // 唯一索引 同mysql唯一索引
		DropDups:   false,              // 索引重复替换旧文档,Unique为true时失效
		Background: true,               // 后台创建索引
	}
}
