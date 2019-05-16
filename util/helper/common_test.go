package helper

import (
	"testing"
)

func TestDistinctStringSlice(t *testing.T) {
	type args struct {
		slice []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestDistinctStringSlice",
			args: args{
				slice: append([]string{"1", "2", "3"}, []string{"2", "4", "6", "3"}...),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := DistinctStringSlice(tt.args.slice)
			t.Log(res)
		})
	}
}
