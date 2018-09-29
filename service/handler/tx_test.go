package handler

import (
	"sync"
	"testing"

	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/helper"
	"os"
)

const (
	BankHeight                   = 17694
	StakeCreateHeight            = 28848
	StakeEditHeight              = 28581
	StakeDelegateHeight          = 79026
	StakeBeginUnbondingHeight    = 79063
	StakeCompleteUnbondingHeight = 79177
)

func TestMain(m *testing.M) {
	// setup
	store.Start()

	code := m.Run()

	// shutdown
	os.Exit(code)
}

func buildDocData(blockHeight int64) document.CommonTx {
	client := helper.GetClient()
	// release client
	defer client.Release()

	block, err := client.Client.Block(&blockHeight)

	if err != nil {
		logger.Error.Panic(err)
	}

	if block.BlockMeta.Header.NumTxs > 0 {
		txs := block.Block.Data.Txs
		txByte := txs[0]
		docTx := helper.ParseTx(codec.Cdc, txByte, block.Block)

		return docTx

	}
	return document.CommonTx{}
}

func TestSaveTx(t *testing.T) {
	type args struct {
		docTx document.CommonTx
		mutex sync.Mutex
	}
	tests := []struct {
		name string
		args args
	}{
		//{
		//	name: "tx bank",
		//	args: args{
		//		docTx: buildDocData(BankHeight),
		//		mutex: sync.Mutex{},
		//	},
		//},
		//{
		//	name: "tx stake/create",
		//	args: args{
		//		docTx: buildDocData(StakeCreateHeight),
		//		mutex: sync.Mutex{},
		//	},
		//},
		//{
		//	name: "tx stake/edit",
		//	args: args{
		//		docTx: buildDocData(StakeEditHeight),
		//		mutex: sync.Mutex{},
		//	},
		//},
		//{
		//	name: "tx stake/delegate",
		//	args: args{
		//		docTx: buildDocData(StakeDelegateHeight),
		//		mutex: sync.Mutex{},
		//	},
		//},
		//{
		//	name: "tx stake/beginUnbonding",
		//	args: args{
		//		docTx: buildDocData(StakeBeginUnbondingHeight),
		//		mutex: sync.Mutex{},
		//	},
		//},
		{
			name: "tx stake/completeUnbonding",
			args: args{
				docTx: buildDocData(StakeCompleteUnbondingHeight),
				mutex: sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SaveTx(tt.args.docTx, tt.args.mutex)
		})
	}
}
