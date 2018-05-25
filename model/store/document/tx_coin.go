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

// Coin tx
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

func (c CoinTx) Index() []mgo.Index {
	return []mgo.Index{
		{
			Key:        []string{"tx_hash"},
			Unique:     true,
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
			Key:        []string{"to"},
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
			Key:        []string{"from", "to"},
			Unique:     false,
			DropDups:   false,
			Background: true,
		},

	}
}
