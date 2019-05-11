// package for parse tx struct from binary data

package helper

import (
	"encoding/json"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//InitClientPool()
	code := m.Run()
	os.Exit(code)
}

func TestParseTx(t *testing.T) {
	client := GetClient()
	// release client
	defer client.Release()

	var height = int64(317356)

	block, err := client.Client.Block(&height)

	if err != nil {
		t.Fatal(err)
	}

	if block.BlockMeta.Header.NumTxs > 0 {
		txs := block.Block.Data.Txs
		tx, accounts := ParseTx(txs[0], block.Block)
		txBytes, _ := json.MarshalIndent(tx, "", "\t")
		t.Logf("tx is %v, accounts is %v\n", string(txBytes), accounts)
	}

}
