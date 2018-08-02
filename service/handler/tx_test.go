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
	//docTxBank := buildDocData(17)
	//docTxStakeCreate := buildDocData(46910)
	docTxStakeBeginUnBond := buildDocData(2753)
	//docTxStakeCompleteUnBond := buildDocData(287)
	//docTxStakeEdit := buildDocData(127)
	//docTxStakeDelegate := buildDocData(81)

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
			name: "tx stake/beginUnbonding",
			args: args{
				docTx: docTxStakeBeginUnBond,
				mutex: sync.Mutex{},
			},
		},
		//{
		//	name: "tx stake/completeUnbonding",
		//	args: args{
		//		docTx: docTxStakeCompleteUnBond,
		//		mutex: sync.Mutex{},
		//	},
		//},
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
		name string
		args args
	}{
		{
			name: "test get validator",
			args: args{
				valAddr: "faa1wp3jgnndsfyxxfeyluu9wsu0yxeseqn6f76fq3",
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
		name string
		args args
	}{
		{
			name: "test get delegation",
			args: args{
				delAddr: "faa1utem9ysq9gkpkhnrrtznmrxyy238kwd0gkcz60",
				valAddr: "faa1wp3jgnndsfyxxfeyluu9wsu0yxeseqn6f76fq3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := getDelegation(tt.args.delAddr, tt.args.valAddr)
			if res.DelegatorAddr == nil {
				logger.Info.Println("delegation is empty")
			}
			if err != nil {
				logger.Error.Fatalln(err)
			}

			logger.Info.Println(helper.ToJson(res))
			logger.Info.Println(res.Shares.RatString())
			logger.Info.Println(res.Shares.Float64())
		})
	}
}
