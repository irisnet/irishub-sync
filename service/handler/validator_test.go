package handler

import (
	"fmt"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/helper"
	"testing"
)

func TestCompareAndUpdateValidators(t *testing.T) {
	c := helper.GetClient()
	defer c.Release()

	status, _ := c.Client.Status()
	res, _ := c.Client.Validators(&status.SyncInfo.LatestBlockHeight)
	tmVals := res.Validators

	type args struct {
		tmVals []*types.Validator
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test compare and update validators",
			args: args{
				tmVals: tmVals,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CompareAndUpdateValidators(tt.args.tmVals)
		})
	}
}

func Test_compareSlice(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test compare two slice: lenght not equal",
			args: args{
				a: nil,
				b: []string{"2", "3"},
			},
		},
		{
			name: "test compare two slice: contain empty",
			args: args{
				a: []string{"1", "2", "3"},
				b: []string{"2", "3", ""},
			},
		},
		{
			name: "test compare two slice: element not equal",
			args: args{
				a: []string{"1", "2", "3"},
				b: []string{"2", "3", "4"},
			},
		},
		{
			name: "test compare two slice: equal",
			args: args{
				a: []string{"1", "2", "3"},
				b: []string{"2", "3", "1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := compareSlice(tt.args.a, tt.args.b)
			fmt.Println(res)
		})
	}
}

func Test_sliceContains(t *testing.T) {
	type args struct {
		s []string
		e string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test slice contain",
			args: args{
				s: []string{"1", "2"},
				e: "1",
			},
		},
		{
			name: "test slice not contain",
			args: args{
				s: []string{"1", "2"},
				e: "1.1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := sliceContains(tt.args.s, tt.args.e)
			fmt.Println(res)
		})
	}
}
