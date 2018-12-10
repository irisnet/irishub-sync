// package for parse tx struct from binary data

package helper

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/types"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//InitClientPool()
	code := m.Run()
	os.Exit(code)
}

func buildTxByte(blockHeight int64) (types.Tx, *types.Block) {
	client := GetClient()
	// release client
	defer client.Release()

	block, err := client.Client.Block(&blockHeight)

	if err != nil {
		logger.Panic(err.Error())
	}

	if block.BlockMeta.Header.NumTxs > 0 {
		txs := block.Block.Data.Txs
		return txs[0], block.Block
	}

	return nil, nil
}
