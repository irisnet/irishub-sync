package document

import (
	"testing"
)

func TestTxGas_RemoveAll(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test remove all data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := TxGas{}
			if err := d.RemoveAll(); err != nil {
				t.Errorf("error = %v\n", err)
			}
		})
	}
}

func TestTxGas_SaveAll(t *testing.T) {
	type args struct {
		txGases []TxGas
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test save all data",
			args: args{
				txGases: []TxGas{
					{
						TxType: "Transfer",
						GasUsed: GasUsed{
							MinGasUsed: 1.0,
							MaxGasUsed: 2.0,
							AvgGasUsed: 1.2,
						},
						GasPrice: GasPrice{
							Denom:       "iris",
							MinGasPrice: 1.1,
							MaxGasPrice: 1.2,
							AvgGasPrice: 1.15,
						},
					},
					{
						TxType: "Delegate",
						GasUsed: GasUsed{
							MinGasUsed: 1.0,
							MaxGasUsed: 2.0,
							AvgGasUsed: 1.2,
						},
						GasPrice: GasPrice{
							Denom:       "iris",
							MinGasPrice: 1.1,
							MaxGasPrice: 1.2,
							AvgGasPrice: 1.15,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := TxGas{}
			if err := d.SaveAll(tt.args.txGases); err != nil {
				t.Errorf("TxGas.SaveAll() error = %v\n", err)
			}
		})
	}
}
