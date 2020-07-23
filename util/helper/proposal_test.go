package helper

import (
	"encoding/json"
	"testing"
)

func TestGetProposal(t *testing.T) {
	proposalId := uint64(1)
	if res, err := GetProposal(proposalId); err != nil {
		t.Fatal(err)
	} else {
		resBytes, _ := json.MarshalIndent(res, "", "\t")
		t.Log(string(resBytes))
	}
}
