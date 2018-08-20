package document

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/util/constant"
	"testing"
)

func TestCommonTx_CalculateTxGasAndGasPrice(t *testing.T) {
	type args struct {
		txType string
		limit  int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test calculate tx gas and gas price",
			args: args{
				txType: constant.TxTypeStakeCreateValidator,
				limit:  10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := CommonTx{}
			res, err := d.CalculateTxGasAndGasPrice(tt.args.txType, tt.args.limit)
			if err != nil {
				t.Error(err)
			} else {
				raw, _ := json.Marshal(res)
				t.Logf("res is %v", string(raw))
			}
		})
	}
}
