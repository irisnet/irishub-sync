package handler

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/helper"
	"sort"
	"testing"
)

func TestCompareAndUpdateValidators(t *testing.T) {
	c := helper.GetClient()
	defer c.Release()

	status, _ := c.Client.Status()
	res, _ := c.Client.Validators(&status.SyncInfo.LatestBlockHeight)
	tmVals := res.Validators

	type args struct {
		tmVals []*types.Validator
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test compare and update validators",
			args: args{
				tmVals: tmVals,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CompareAndUpdateValidators()
		})
	}
}

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s: %d", p.Name, p.Age)
}

func TestSort(t *testing.T) {
	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	sort.SliceStable(people, func(i, j int) bool {
		return people[i].Age > people[j].Age
	})
	fmt.Println(people)
}

func TestBuildUnbondingDelegation(t *testing.T) {
	var delAddr = "faa1ljemm0yznz58qxxs8xyak7fashcfxf5lssn6jm"
	var valAddr = "fva1kca5vw7r2k72d5zy0demszmrhdz4dp8t4uat0c"
	res := BuildUnbondingDelegation(delAddr, valAddr)
	r, _ := json.Marshal(res)
	fmt.Println(string(r))
}

func TestBuildDelegation(t *testing.T) {
	var delAddr = "faa1ljemm0yznz58qxxs8xyak7fashcfxf5lssn6jm"
	var valAddr = "fva1kca5vw7r2k72d5zy0demszmrhdz4dp8t4uat0c"
	res := BuildDelegation(delAddr, valAddr)
	r, _ := json.Marshal(res)
	fmt.Println(string(r))
}
