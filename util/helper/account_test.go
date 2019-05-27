// This package is used for Query balance of account

package helper

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"encoding/hex"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/stretchr/testify/assert"
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

func TestConvertAccountAddrFromHexToBech32(t *testing.T) {
	hexAddr := "7c99a0d9ab962250b83234f830666d785e5406ff"
	bech32Addr := "faa10jv6pkdtjc39pwpjxnurqend0p09gphl3xg5yc"
	if bytes, err := hex.DecodeString(hexAddr); err != nil {
		t.Fatal(err)
	} else {
		res, err := ConvertAccountAddrFromHexToBech32(bytes)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, bech32Addr, res)
	}

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
