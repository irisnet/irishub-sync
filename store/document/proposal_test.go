package document

import (
	"testing"
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-sync/store"
)

func TestQueryProposal(t *testing.T) {
	store.Start()
	data, err := QueryProposal(10)
	if err != nil {
		t.Fatal(err)
	}
	data.Description = "Proposal Test"
	err = store.SaveOrUpdate(data)
	fmt.Println(err)
	strval, _ := json.Marshal(data)
	fmt.Println(string(strval))
	store.Stop()
}
