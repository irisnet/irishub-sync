package document

import (
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmStakeTx = "tx_stake"
)

// StakeTx
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
	return CollectionNmStakeTx
}

func (c StakeTx) PkKvPair() map[string]interface{} {
	return bson.M{"tx_hash": c.TxHash}
}

func (c StakeTx) Index() []mgo.Index {
	return []mgo.Index{
		{
			Key:        []string{"tx_hash"},
			Unique:     false,
			DropDups:   false,
			Background: true,
		},
		{
			Key:        []string{"from"},
			Unique:     false,
			DropDups:   false,
			Background: true,
		},
		{
			Key:        []string{"pub_key"},
			Unique:     false,
			DropDups:   false,
			Background: true,
		},
		{
			Key:        []string{"-height"},
			Unique:     false,
			DropDups:   false,
			Background: true,
		},
		{
			Key:        []string{"type"},
			Unique:     false,
			DropDups:   false,
			Background: true,
		},
		{
			Key:        []string{"from", "pub_key", "type"},
			Unique:     false,
			DropDups:   false,
			Background: true,
		},
	}
}
