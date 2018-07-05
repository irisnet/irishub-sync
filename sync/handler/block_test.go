package handler

import (
	"testing"

	"github.com/tendermint/tendermint/types"
	"github.com/irisnet/irishub-sync/util/helper"
	"github.com/irisnet/irishub-sync/module/logger"
)

func buildBlock(blockHeight int64) (*types.BlockMeta, *types.Block) {

	client := helper.GetClient()
	// release client
	defer client.Release()

	block, err := client.Client.Block(&blockHeight)

	if err != nil {
		logger.Error.Fatalln(err)
	}

	return block.BlockMeta, block.Block
}

func TestSaveBlock(t *testing.T) {
	meta1, block1 := buildBlock(28558)
	meta2, block2 := buildBlock(47349)

	type args struct {
		meta  *types.BlockMeta
		block *types.Block
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test save block",
			args: args{
				meta: meta1,
				block: block1,
			},
		},
		{
			name: "test save block",
			args: args{
				meta: meta2,
				block: block2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SaveBlock(tt.args.meta, tt.args.block)
		})
	}
}
