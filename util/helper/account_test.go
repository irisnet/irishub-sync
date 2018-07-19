// This package is used for query balance of account

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
				address: "D770D45DEA7548076F8A27F9C9749B200934F1B4",
			},
		},
		{
			name: "test balance is nil",
			args: args{
				address: "D770D45DEA7548076F8A27F9C9749B200934F1B9",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QueryAccountBalance(tt.args.address)
			logger.Info.Println(ToJson(got))
		})
	}
}
