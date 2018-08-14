package handler

import (
	"sync"
	"testing"

	"github.com/irisnet/irishub-sync/module/codec"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/helper"
)

const (
	BankHeight                   = 17694
	StakeCreateHeight            = 28848
	StakeEditHeight              = 28581
	StakeDelegateHeight          = 29209
	StakeBeginUnbondingHeight    = 29214
	StakeCompleteUnbondingHeight = 29306
)

func init() {
	helper.InitClientPool()
	store.InitWithAuth()
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
		{
			name: "tx bank",
			args: args{
				docTx: buildDocData(BankHeight),
				mutex: sync.Mutex{},
			},
		},
		{
			name: "tx stake/create",
			args: args{
				docTx: buildDocData(StakeCreateHeight),
				mutex: sync.Mutex{},
			},
		},
		{
			name: "tx stake/edit",
			args: args{
				docTx: buildDocData(StakeEditHeight),
				mutex: sync.Mutex{},
			},
		},
		{
			name: "tx stake/delegate",
			args: args{
				docTx: buildDocData(StakeDelegateHeight),
				mutex: sync.Mutex{},
			},
		},
		{
			name: "tx stake/beginUnbonding",
			args: args{
				docTx: buildDocData(StakeBeginUnbondingHeight),
				mutex: sync.Mutex{},
			},
		},
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
				valAddr: "faa15lpdxlk0hwkewmncdhlyfle8jc3k9xzhh75txs",
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
