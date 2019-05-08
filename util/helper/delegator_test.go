package helper

import (
	"encoding/json"
	"testing"
)

func TestGetDelegations(t *testing.T) {
	type args struct {
		delegator string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestGetDelegations",
			args: args{
				delegator: "faa1eqvkfthtrr93g4p9qspp54w6dtjtrn279vcmpn",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := GetDelegations(tt.args.delegator)
			resBytes, _ := json.MarshalIndent(res, "", "\t")
			t.Log(string(resBytes))
		})
	}
}

func TestGetUnbondingDelegations(t *testing.T) {
	type args struct {
		delegator string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestGetDelegations",
			args: args{
				delegator: "faa1eqvkfthtrr93g4p9qspp54w6dtjtrn279vcmpn",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := GetUnbondingDelegations(tt.args.delegator)
			resBytes, _ := json.MarshalIndent(res, "", "\t")
			t.Log(string(resBytes))
		})
	}
}
