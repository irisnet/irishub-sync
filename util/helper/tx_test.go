// package for parse tx struct from binary data

package helper

import (
	"testing"

	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/tendermint/tendermint/types"
	"github.com/irisnet/irishub-sync/module/logger"
)

func buildTxByte(blockHeight int64) (types.Tx, *types.Block) {

	InitClientPool()
	client := GetClient()
	// release client
	defer client.Release()

	block, err := client.Client.Block(&blockHeight)

	if err != nil {
		logger.Error.Panic(err)
	}

	if block.BlockMeta.Header.NumTxs > 0 {
		txs := block.Block.Data.Txs
		return txs[0], block.Block
	}

	return nil, nil
}

func TestParseTx(t *testing.T) {
	coinByte, coinBlock := buildTxByte(1762)
	scByte, scBlock := buildTxByte(46910)
	seByte, seBlock := buildTxByte(49388)
	sdByte, sdBlock := buildTxByte(47349)
	suByte, suBlock := buildTxByte(34241)

	type args struct {
		txByte types.Tx
		block *types.Block
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 interface{}
	}{
		{
			name: "test tx coin",
			args: args{
				txByte: coinByte,
				block: coinBlock,
			},
		},
		{
			name: "test tx stake/create",
			args: args{
				txByte: scByte,
				block: scBlock,
			},
		},
		{
			name: "test tx stake/edit",
			args: args{
				txByte: seByte,
				block: seBlock,
			},
		},
		{
			name: "test tx stake/delegate",
			args: args{
				txByte: sdByte,
				block: sdBlock,
			},
		},
		{
			name: "test tx stake/unbond",
			args: args{
				txByte: suByte,
				block: suBlock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ParseTx(codec.Cdc, tt.args.txByte, tt.args.block)
			logger.Info.Printf("Tx is %v\n", ToJson(res))
		})
	}
}
