package document

import (
	"encoding/json"
	"testing"
)

func TestGetAccountByPK(t *testing.T) {
	doc := Account{
		Address: "123",
	}
	res, err := doc.getAccountByPK()
	if err != nil {
		t.Fatal(err)
	}
	resBytes, _ := json.MarshalIndent(res, "", "\t")
	t.Log(string(resBytes))
}


