package handler

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

func HandleTx(block *types.Block) error {
	var (
		batch []txn.Op
	)
	for _, txByte := range block.Txs {
		tx := helper.ParseTx(txByte, block)
		txOp := txn.Op{
			C:      document.CollectionNmCommonTx,
			Id:     bson.NewObjectId(),
			Insert: tx,
		}
		batch = append(batch, txOp)

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
		SaveOrUpdateAccountDelegationInfo(tx)

		switch tx.Type {
		case constant.TxTypeStakeBeginUnbonding:
			accounts := []string{tx.From}
			SaveOrUpdateAccountUnbondingDelegationInfo(accounts, tx.Height, tx.Time.Unix())
		}

	}

	if len(batch) > 0 {
		err := store.Txn(batch)
		if err != nil {
			return err
		}
	}

	return nil
}
