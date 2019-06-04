package handler

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestIsContainVotingPeriodStartTag(t *testing.T) {
	txHash := "7DD721FAFF970A6A74D1BF0771A1634B10CF40CCDD10C8F1A38F65BB0E035B2D"

	var tx document.CommonTx
	fn := func(c *mgo.Collection) error {
		q := bson.M{"tx_hash": txHash}
		return c.Find(q).One(&tx)
	}

	if err := store.ExecCollection(tx.Name(), fn); err != nil {
		t.Fatal(err)
	} else {
		res := isContainVotingPeriodStartTag(tx)
		t.Log(res)
	}
}

func TestHandleProposal(t *testing.T) {
	txHash := "2D9B3B49F6250B8A2D60C90AF0F21591CDDC74C14FBF7E1ADDE0062D1E977922"

	var tx document.CommonTx
	fn := func(c *mgo.Collection) error {
		q := bson.M{"tx_hash": txHash}
		return c.Find(q).One(&tx)
	}

	if err := store.ExecCollection(tx.Name(), fn); err != nil {
		t.Fatal(err)
	} else {
		handleProposal(tx)
		t.Log("success")
	}
}
