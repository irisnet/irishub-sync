// package for parse tx struct from binary data

package helper

import (
	"fmt"
	"github.com/irisnet/irishub-sync/logger"
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

	var height = int64(103)

	block, err := client.Client.Block(&height)

	if err != nil {
		logger.Panic(err.Error())
	}

	if block.BlockMeta.Header.NumTxs > 0 {
		txs := block.Block.Data.Txs
		tx := ParseTx(txs[0], block.Block)
		fmt.Println(tx)
	}

}

func TestParseTx(t *testing.T) {
	client := GetClient()
	// release client
	defer client.Release()

	height := int64(52373)
	block, err := client.Block(&height)

	if err != nil {
		logger.Panic(err.Error())
	}

	if block.BlockMeta.Header.NumTxs > 0 {
		txs := block.Block.Data.Txs
		for _, tx := range txs {
			ParseTx(tx, block.Block)
		}
	}

}
