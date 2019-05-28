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
				delegator: "faa192vef4442d07lqde59mx35dvmfv9v72wrsu84a",
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
				delegator: "faa192vef4442d07lqde59mx35dvmfv9v72wrsu84a",
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

func TestCalculateDelegatorDelegationTokens(t *testing.T) {
	delegatorAddr := "faa192vef4442d07lqde59mx35dvmfv9v72wrsu84a"
	tokens := CalculateDelegatorDelegationTokens(GetDelegations(delegatorAddr))
	t.Log(tokens)
}

func TestCalculateDelegatorUnbondingDelegationTokens(t *testing.T) {
	delegatorAddr := "faa192vef4442d07lqde59mx35dvmfv9v72wrsu84a"
	tokens := CalculateDelegatorUnbondingDelegationTokens(GetUnbondingDelegations(delegatorAddr))
	t.Log(tokens)
}
