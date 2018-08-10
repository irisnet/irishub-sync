package handler

import (
	"sync"
	"testing"

	"github.com/irisnet/irishub-sync/store"
)

func TestSaveAccount(t *testing.T) {
	docTxBank := buildDocData(17)
	//docTxStakeCreate := buildDocData(46910)
	docTxStakeBeginUnBonding := buildDocData(148)
	docTxStakeCompleteUnBonding := buildDocData(287)
	docTxStakeEdit := buildDocData(127)
	docTxStakeDelegate := buildDocData(81)

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
			name: "tx stake/beginUnbonding",
			args: args{
				docTx: docTxStakeBeginUnBonding,
				mutex: sync.Mutex{},
			},
		},
		{
			name: "tx stake/completeUnbonding",
			args: args{
				docTx: docTxStakeCompleteUnBonding,
				mutex: sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SaveAccount(tt.args.docTx, tt.args.mutex)
		})
	}
}

func TestUpdateBalance(t *testing.T) {
	docTxBank := buildDocData(17)
	//docTxStakeCreate := buildDocData(46910)
	docTxStakeBeginUnBonding := buildDocData(148)
	docTxStakeCompleteUnBonding := buildDocData(287)
	docTxStakeEdit := buildDocData(127)
	docTxStakeDelegate := buildDocData(81)

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
			name: "tx stake/beginUnbonding",
			args: args{
				docTx: docTxStakeBeginUnBonding,
				mutex: sync.Mutex{},
			},
		},
		{
			name: "tx stake/completeUnbonding",
			args: args{
				docTx: docTxStakeCompleteUnBonding,
				mutex: sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateBalance(tt.args.docTx, tt.args.mutex)
		})
	}
}
