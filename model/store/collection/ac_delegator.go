package collection

import (
	"errors"
	"github.com/irisnet/iris-sync-server/model/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	DocsNmDelegator = "ac_delegator"
)

type Delegator struct {
	Address string `bson:"address"`
	PubKey  string `bson:"pub_key"`
	Shares  int64  `bson:"shares"`
}

func (d Delegator) Name() string {
	return DocsNmDelegator
}

func (d Delegator) PkKvPair() map[string]interface{} {
	return bson.M{"address": d.Address, "pub_key": d.PubKey}
}

func (d Delegator) Index() mgo.Index {
	return mgo.Index{
		Key:        []string{"address"}, // 索引字段， 默认升序,若需降序在字段前加-
		Unique:     false,               // 唯一索引 同mysql唯一索引
		DropDups:   false,               // 索引重复替换旧文档,Unique为true时失效
		Background: true,                // 后台创建索引
	}
}

func QueryDelegatorByAddressAndPubkey(address string, pubKey string) (Delegator, error) {
	var result Delegator
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{"address": address, "pub_key": pubKey}).Sort("-shares").One(&result)
		return err
	}

	if store.ExecCollection(DocsNmDelegator, query) != nil {
		log.Printf("delegator is Empty")
		return result, errors.New("delegator is Empty")
	}

	return result, nil
}
