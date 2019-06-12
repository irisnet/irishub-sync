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
	txHash := "A837138C2A569B7884AA94C27CC4AB791C04F1B8DD93EFC3D5BFCF3D7EB0F2F3"

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
	txHash := "5875062EE8B8656CF943C42983F382B5341B1B0C530062D266BD8283CA9658B0"

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
