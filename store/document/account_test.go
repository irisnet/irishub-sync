package document

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/constant"
	"testing"
	"time"
)

func TestUpsert(t *testing.T) {
	account := Account{
		Address:       "faa1eqvkfthtrr93g4p9qspp54w6dtjtrn279vcmpn",
		AccountNumber: uint64(12),
		CoinIris: store.Coin{
			Denom:  constant.IrisAttoUnit,
			Amount: float64(123456),
		},
		CoinIrisUpdateAt: time.Now().Unix(),
	}
	if err := store.Upsert(account); err != nil {
		t.Fatal(err)
	}
}
