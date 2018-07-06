package handler

import (
	"sync"
	"testing"

	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/helper"
)

func init() {
	helper.InitClientPool()
	store.Init()
}

func buildDocData(blockHeight int64) store.Docs {

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
	return nil
}

func TestSaveTx(t *testing.T) {
	//docTxBank := buildDocData(1762)
	//docTxStakeCreate := buildDocData(46910)
	//docTxStakeEdit := buildDocData(49388)
	//docTxStakeDelegate := buildDocData(47349)
	docTxStakeUnBond := buildDocData(96319)

	type args struct {
		docTx store.Docs
		mutex sync.Mutex
	}
	tests := []struct {
		name string
		args args
	}{
		//{
		//	name: "tx bank",
		//	args: args{
		//		docTx: docTxBank,
		//		mutex: sync.Mutex{},
		//	},
		//},
		//{
		//	name: "tx stake/create",
		//	args: args{
		//		docTx: docTxStakeCreate,
		//		mutex: sync.Mutex{},
		//	},
		//},
		//{
		//	name: "tx stake/edit",
		//	args: args{
		//		docTx: docTxStakeEdit,
		//		mutex: sync.Mutex{},
		//	},
		//},
		//{
		//	name: "tx stake/delegate",
		//	args: args{
		//		docTx: docTxStakeDelegate,
		//		mutex: sync.Mutex{},
		//	},
		//},
		{
			name: "tx stake/unbond",
			args: args{
				docTx: docTxStakeUnBond,
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
