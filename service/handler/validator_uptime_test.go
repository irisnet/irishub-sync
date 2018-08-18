package handler

import "testing"

func TestCalculateAndSaveValidatorUptime(t *testing.T) {
	type args struct {
		latestHeight int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test calculate and save validator uptime",
			args: args{
				latestHeight: 90106,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CalculateAndSaveValidatorUptime(tt.args.latestHeight)
		})
	}
}
