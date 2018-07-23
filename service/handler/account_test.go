package handler

import (
	"sync"
	"testing"

	"github.com/irisnet/irishub-sync/store"
)

func TestSaveAccount(t *testing.T) {
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
			SaveAccount(tt.args.docTx, tt.args.mutex)
		})
	}
}

func TestUpdateBalance(t *testing.T) {
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
			UpdateBalance(tt.args.docTx, tt.args.mutex)
		})
	}
}
