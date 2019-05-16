package helper

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetUnbondingDelegation(t *testing.T) {
	var delAddr = "faa1ljemm0yznz58qxxs8xyak7fashcfxf5lssn6jm"
	var valAddr = "fva1kca5vw7r2k72d5zy0demszmrhdz4dp8t4uat0c"

	res := GetUnbondingDelegation(delAddr, valAddr)
	r, _ := json.Marshal(res)
	fmt.Println(string(r))
}

func TestGetValidator(t *testing.T) {
	var valAddr = "fva1phst8wkk27jd748p0nmffzh6288kldlpxq39h8"
	if res, err := GetValidator(valAddr); err != nil {
		t.Fatal(err)
	} else {
		resBytes, _ := json.MarshalIndent(res, "", "\t")
		t.Log(string(resBytes))
	}
}
