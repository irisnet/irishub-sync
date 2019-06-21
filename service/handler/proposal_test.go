package handler

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	itypes "github.com/irisnet/irishub-sync/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestIsContainVotingPeriodStartTag(t *testing.T) {
	txHash := "37A0127A87AA68BFE73D03C2B9A2A6A3D8E51DF242D86C845DB2D158B1617502"

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
	txHash := "37A0127A87AA68BFE73D03C2B9A2A6A3D8E51DF242D86C845DB2D158B1617502"

	var tx document.CommonTx
	fn := func(c *mgo.Collection) error {
		q := bson.M{"tx_hash": txHash}
		return c.Find(q).One(&tx)
	}

	if err := store.ExecCollection(tx.Name(), fn); err != nil {
		t.Fatal(err)
	} else {
		var txMsg document.TxMsg
		fn := func(c *mgo.Collection) error {
			q := bson.M{"hash": txHash}
			return c.Find(q).One(&txMsg)
		}
		if err := store.ExecCollection(txMsg.Name(), fn); err != nil {
			t.Fatal(err)
		} else {
			var msgVote itypes.Vote
			if err := json.Unmarshal([]byte(txMsg.Content), &msgVote); err != nil {
				t.Fatal(err)
			} else {
				tx.Msg = msgVote
				handleProposal(tx)
				t.Log("success")
			}
		}

	}
}
