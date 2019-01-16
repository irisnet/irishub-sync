package types

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewSubmitSoftwareUpgradeProposal(t *testing.T) {
	submitProposal := MsgSubmitProposal{
		Title:       "title",
		Description: "description",
	}
	msg := MsgSubmitSoftwareUpgradeProposal{
		MsgSubmitProposal: submitProposal,
		Version:           uint64(1),
		Software:          "aa",
		SwitchHeight:      uint64(2),
	}

	bz, _ := json.Marshal(msg)
	fmt.Println(string(bz))

	var p MsgSubmitProposal
	json.Unmarshal(bz, &p)

	bz1, _ := json.Marshal(p)
	fmt.Println(string(bz1))
}
