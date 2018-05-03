package document

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmBlock = "block"
)

type Block struct {
	Height int64     `bson:"height"`
	Time   time.Time `bson:"time"`
	TxNum  int64     `bson:"tx_num"`
}

func (d Block) Name() string {
	return CollectionNmBlock
}

func (d Block) PkKvPair() map[string]interface{} {
	return bson.M{"height": d.Height}
}

func (d Block) Index() []mgo.Index {
	return []mgo.Index{
		{
			Key:        []string{"height"},
			Unique:     true,
			DropDups:   false,
			Background: true,
		},
	}
}
