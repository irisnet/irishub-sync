// This package is used for Query balance of account

package helper

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/irisnet/irishub-sync/logger"
)

func TestQueryAccountBalance(t *testing.T) {
	//InitClientPool()

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
				address: "faa1eqvkfthtrr93g4p9qspp54w6dtjtrn279vcmpn",
			},
		},
		//{
		//	name: "test balance is nil",
		//	args: args{
		//		address: "faa1utem9ysq9gkpkhnrrtznmrxyy238kwd0gkcz60",
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, accNumber := QueryAccountInfo(tt.args.address)
			logger.Info("accNum info", logger.Uint64("accNumber", accNumber))
			logger.Info(ToJson(res))
		})
	}
}

func TestValAddrToAccAddr(t *testing.T) {
	valAddr := "fva1qz47703lujvyumg4k3fgl7uf9v7uruhzqqh5f8"
	fmt.Println(ValAddrToAccAddr(valAddr))
}

type Student struct {
	Name   string   `json:"name"`
	Age    int      `json:"age"`
	Course []Course `json:"course"`
}

type Course struct {
	Name     string     `json:"name"`
	Schedule []Schedule `json:"schedule"`
}

type Schedule struct {
	Time time.Time
}

func TestMap2Struct(t *testing.T) {
	data := Student{
		Name: "zhansan",
		Age:  10,
		Course: []Course{
			{Name: "MAth", Schedule: []Schedule{
				{Time: time.Now().UTC()},
			}},
		},
	}

	mp := Struct2Map(data)

	var data1 Student
	Map2Struct(mp, &data1)

	require.EqualValues(t, data, data1)

}
