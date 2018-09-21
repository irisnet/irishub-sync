package handler

import (
	"testing"
)

func TestCalculateTxGasAndGasPrice(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test calculate tx gas and gas price",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CalculateTxGasAndGasPrice()
		})
	}
}
