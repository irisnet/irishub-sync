package handler

import (
	"sync"
	"testing"

	"github.com/irisnet/irishub-sync/store/document"
)

func TestSaveAccount(t *testing.T) {

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
			SaveAccount(tt.args.docTx, tt.args.mutex)
		})
	}
}

func TestUpdateBalance(t *testing.T) {
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
			UpdateBalance(tt.args.docTx, tt.args.mutex)
		})
	}
}
