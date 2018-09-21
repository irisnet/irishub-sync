// This package is used for Query balance of account

package helper

import (
	"testing"

	"github.com/irisnet/irishub-sync/module/logger"
)

func TestQueryAccountBalance(t *testing.T) {
	InitClientPool()

	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test balance not nil",
			args: args{
				address: "faa1r0ljqhd7vwrpwh8h8fa5luh89nljrnkqcdgfr0",
			},
		},
		//{
		//	name: "test balance is nil",
		//	args: args{
		//		address: "faa1utem9ysq9gkpkhnrrtznmrxyy238kwd0gkcz60",
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QueryAccountBalance(tt.args.address)
			logger.Info.Println(ToJson(got))
		})
	}
}
