package document

import (
	"testing"
)

func TestValidatorUpTime_RemoveAll(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test remove all data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := ValidatorUpTime{}
			err := d.RemoveAll()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestValidatorUpTime_SaveAll(t *testing.T) {
	type args struct {
		validatorUpTimes []ValidatorUpTime
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test save all data",
			args: args{
				validatorUpTimes: []ValidatorUpTime{
					{
						ValAddress: "1",
						UpTime:     98.3,
					},
					{
						ValAddress: "2",
						UpTime:     96.3,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := ValidatorUpTime{}
			if err := d.SaveAll(tt.args.validatorUpTimes); err != nil {
				t.Errorf("err is %v\n", err)
			}
		})
	}
}
