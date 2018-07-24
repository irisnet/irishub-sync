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
	store.InitWithAuth()
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
	docTxBank := buildDocData(1707)
	//docTxStakeCreate := buildDocData(46910)
	docTxStakeUnBond := buildDocData(5240)
	docTxStakeEdit := buildDocData(17026)
	docTxStakeDelegate := buildDocData(1760)

	type args struct {
		docTx store.Docs
		mutex sync.Mutex
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "tx bank",
			args: args{
				docTx: docTxBank,
				mutex: sync.Mutex{},
			},
		},
		//{
		//	name: "tx stake/create",
		//	args: args{
		//		docTx: docTxStakeCreate,
		//		mutex: sync.Mutex{},
		//	},
		//},
		{
			name: "tx stake/edit",
			args: args{
				docTx: docTxStakeEdit,
				mutex: sync.Mutex{},
			},
		},
		{
			name: "tx stake/delegate",
			args: args{
				docTx: docTxStakeDelegate,
				mutex: sync.Mutex{},
			},
		},
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

func Test_getValidator(t *testing.T) {
	type args struct {
		valAddr string
	}
	tests := []struct {
		name    string
		args    args
	}{
		{
			name: "test get validator",
			args: args{
				valAddr: "441EF0233B416CF486E21D0377B8758F25FECEAB",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := getValidator(tt.args.valAddr)
			if err != nil {
				logger.Error.Fatalln(err)
			}
			logger.Info.Println(helper.ToJson(res))
		})
	}
}

func Test_getDelegation(t *testing.T) {
	type args struct {
		delAddr string
		valAddr string
	}
	tests := []struct {
		name    string
		args    args
	}{
		{
			name: "test get delegation",
			args: args{
				delAddr: "7E2E6D4764016042B8A82B6EC8B041C79FE5580C",
				valAddr: "441EF0233B416CF486E21D0377B8758F25FECEAB",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := getDelegation(tt.args.delAddr, tt.args.valAddr)
			if err != nil {
				logger.Error.Fatalln(err)
			}
			logger.Info.Println(helper.ToJson(res))
		})
	}
}
