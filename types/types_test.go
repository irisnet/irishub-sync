package types

import (
	"encoding/json"
	"testing"
)

func TestParseCoin(t *testing.T) {
	coin := ParseCoin("11111lc1")
	coinStr, _ := json.Marshal(coin)
	t.Log(string(coinStr))
}
