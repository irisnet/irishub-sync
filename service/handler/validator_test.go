package handler

import (
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
			CompareAndUpdateValidators()
		})
	}
}
