package handler

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/helper"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

func HandleTx(block *types.Block) (error) {
	var (
		batch                  []txn.Op
	)

	for _, txByte := range block.Txs {
		tx := helper.ParseTx(txByte, block)

		// batch insert tx
		txOp := txn.Op{
			C:      document.CollectionNmCommonTx,
			Id:     bson.NewObjectId(),
			Insert: tx,
		}
		batch = append(batch, txOp)

		// save or update proposal
		handleProposal(tx)
		//handleTokenFlow(blockWithTags, tx, &batch)

		// save new account address
		saveNewAccount(&tx)
	}

	if len(batch) > 0 {
		err := store.Txn(batch)
		if err != nil {
			return err
		}
	}

	return nil
}
