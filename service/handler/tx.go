package handler

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/helper"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

func HandleTx(block *types.Block) error {
	var (
		batch           []txn.Op
		accountsInBlock []string
	)
	for _, txByte := range block.Txs {
		tx, accountsInTx := helper.ParseTx(txByte, block)
		txOp := txn.Op{
			C:      document.CollectionNmCommonTx,
			Id:     bson.NewObjectId(),
			Insert: tx,
		}
		batch = append(batch, txOp)
		accountsInBlock = append(accountsInBlock, accountsInTx...)

		msg := tx.Msg
		if msg != nil {
			txMsg := document.TxMsg{
				Hash:    tx.TxHash,
				Type:    msg.Type(),
				Content: msg.String(),
			}
			txOp := txn.Op{
				C:      document.CollectionNmTxMsg,
				Id:     bson.NewObjectId(),
				Insert: txMsg,
			}
			batch = append(batch, txOp)
		}
		// TODO(deal with by biz system)
		handleProposal(tx)
		SaveOrUpdateDelegator(tx)
	}

	if len(batch) > 0 {
		err := store.Txn(batch)
		if err != nil {
			return err
		}
	}

	// update account balance
	// don't use goroutine for this method, sync already use multiple goroutine to execute task,
	// if task goroutine contain other goroutine, the number of goroutine will out of control
	UpdateAccountInfo(accountsInBlock, block.Time.Unix())

	return nil
}
