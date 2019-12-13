package handler

import (
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"strings"
)

func HandleTx(block *types.Block) ([]string, error) {
	var (
		batch                  []txn.Op
		accsBalanceNeedUpdated []string
	)
	getAccsBalanceNeedUpdated := func(addr string) {
		if strings.HasPrefix(addr, types.Bech32AccountAddrPrefix) {
			accsBalanceNeedUpdated = append(accsBalanceNeedUpdated, addr)
		}
	}

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

		// save or update account delegations info and unbonding delegation info
		SaveOrUpdateAccountDelegationInfo(tx)
		switch tx.Type {
		case constant.TxTypeStakeBeginUnbonding:
			accounts := []string{tx.From}
			SaveOrUpdateAccountUnbondingDelegationInfo(accounts, tx.Height, tx.Time.Unix())
		}

		// get accounts which balance need updated by parse tx
		getAccsBalanceNeedUpdated(tx.From)
		getAccsBalanceNeedUpdated(tx.To)
	}

	if len(batch) > 0 {
		err := store.Txn(batch)
		if err != nil {
			return accsBalanceNeedUpdated, err
		}
	}

	return accsBalanceNeedUpdated, nil
}
