package handler

import (
	"testing"

	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/util/helper"
	"github.com/tendermint/tendermint/types"
)

func buildBlock(blockHeight int64) (*types.BlockMeta, *types.Block, []*types.Validator) {

	client := helper.GetClient()
	// release client
	defer client.Release()

	block, err := client.Client.Block(&blockHeight)

	if err != nil {
		logger.Error.Fatalln(err)
	}

	validators, err := client.Client.Validators(&blockHeight)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	return block.BlockMeta, block.Block, validators.Validators
}

func TestSaveBlock(t *testing.T) {
	meta1, block1, vals1 := buildBlock(516)
	meta2, block2, vals2 := buildBlock(515)
	//meta3, block3, vals3 := buildBlock(34241)

	type args struct {
		meta  *types.BlockMeta
		block *types.Block
		vals  []*types.Validator
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test save block",
			args: args{
				meta:  meta1,
				block: block1,
				vals:  vals1,
			},
		},
		{
			name: "test save block",
			args: args{
				meta:  meta2,
				block: block2,
				vals:  vals2,
			},
		},
		//{
		//	name: "test save block",
		//	args: args{
		//		meta:  meta3,
		//		block: block3,
		//		vals:  vals3,
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SaveBlock(tt.args.meta, tt.args.block, tt.args.vals)
		})
	}
}
