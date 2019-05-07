package document

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/constant"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

func TestUpsert(t *testing.T) {
	account := Account{
		Address:       "faa1eqvkfthtrr93g4p9qspp54w6dtjtrn279vcmpn2",
		AccountNumber: 11,
	}
	update := bson.M{
		"$set": bson.M{
			"coin_iris": bson.M{
				"denom":  constant.IrisAttoUnit,
				"amount": 1236,
			},
			"coin_iris_update_at": time.Now().Unix(),
		},
	}
	if err := store.Upsert(account, update); err != nil {
		t.Fatal(err)
	}
}

func TestSave(t *testing.T) {
	account := Account{
		Address:       "faa1eqvkfthtrr93g4p9qspp54w6dtjtrn279vcmpn",
		AccountNumber: uint64(12),
		CoinIris: store.Coin{
			Denom:  constant.IrisAttoUnit,
			Amount: float64(123456),
		},
		CoinIrisUpdateAt: time.Now().Unix(),
	}

	if err := store.Save(account); err != nil {
		t.Fatal(err)
	}
	t.Log("success")
}
