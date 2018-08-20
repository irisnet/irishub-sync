package handler

import "testing"

func TestCalculateAndSaveValidatorUptime(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			name: "test calculate and save validator uptime",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CalculateAndSaveValidatorUpTime()
		})
	}
}
