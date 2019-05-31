package handler

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/helper"
)

func buildBlock(blockHeight int64) (*types.BlockMeta, *types.Block, []*types.Validator) {

	client := helper.GetClient()
	// release client
	defer client.Release()

	block, err := client.Client.Block(&blockHeight)

	if err != nil {
		logger.Error(err.Error())
	}

	validators, err := client.Client.Validators(&blockHeight)
	if err != nil {
		logger.Error(err.Error())
	}

	return block.BlockMeta, block.Block, validators.Validators
}

func TestForEach(t *testing.T) {
	var i int
	var arr = []string{"1", "2", "3"}
	for i = range arr {
		fmt.Println(fmt.Sprintf("a[%d] = %s", i, arr[i]))
		if arr[i] == "2" {
			break
		}

	}
	fmt.Println(fmt.Sprintf("a[%d]", i))
}

func TestParseBlockResult(t *testing.T) {
	v := parseBlockResult(213637)
	bz, _ := json.Marshal(v)
	fmt.Println(string(bz))
}

func TestParseBlock(t *testing.T) {
	blockHeight := int64(88)

	client := helper.GetClient()
	defer client.Release()

	if res, err := client.Block(&blockHeight); err != nil {
		t.Fatal(err)
	} else {
		var validators []*types.Validator
		valRes, err := client.Validators(&blockHeight)
		if err != nil {
			t.Error(err)
		} else {
			validators = valRes.Validators
		}
		blockDoc := ParseBlock(res.BlockMeta, res.Block, validators, nil)

		resBytes, _ := json.MarshalIndent(blockDoc, "", "\t")
		t.Log(string(resBytes))
	}
}
