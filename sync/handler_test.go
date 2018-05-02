package sync

import (
	"testing"

	"github.com/irisnet/iris-sync-server/model/store"

	"github.com/irisnet/iris-sync-server/util/helper"
	"github.com/irisnet/iris-sync-server/module/logger"
	"github.com/irisnet/iris-sync-server/util/constant"
	"github.com/irisnet/iris-sync-server/model/store/collection"
	"strings"
	"encoding/hex"
)

func buildDocData(blockHeight int64) store.Docs {

	helper.InitClientPool()
	client := helper.GetClient()
	// release client
	defer client.Release()

	block, err := client.Client.Block(&blockHeight)

	if err != nil {
		logger.Error.Panic(err)
	}

	if block.BlockMeta.Header.NumTxs > 0 {
		txs := block.Block.Data.Txs
		txByte := txs[0]
		txType, tx := helper.ParseTx(txByte)
		txHash := strings.ToUpper(hex.EncodeToString(txByte.Hash()))

		switch txType {
		case constant.TxTypeCoin:
			coinTx, _ := tx.(document.CoinTx)
			coinTx.TxHash = txHash
			coinTx.Height = block.Block.Height
			coinTx.Time = block.Block.Time
			return coinTx
		case constant.TxTypeStake:
			stakeTx, _ := tx.(document.StakeTx)
			stakeTx.TxHash = txHash
			stakeTx.Height = block.Block.Height
			stakeTx.Time = block.Block.Time
			return stakeTx
		}
	}
	return nil
}

func Test_saveTx(t *testing.T) {
	doc_tx_coin := buildDocData(12453)
	doc_tx_stake_declareCandidacy := buildDocData(19073)
	doc_tx_stake_delegate := buildDocData(13725)
	doc_tx_stake_unbond := buildDocData(14260)

	type args struct {
		tx store.Docs
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name:"save tx_coin",
			args: struct{ tx store.Docs }{
				tx: doc_tx_coin,},
		},
		{
			name:"save tx_stake_declareCandidacy",
			args: struct{ tx store.Docs }{
				tx: doc_tx_stake_declareCandidacy,},

		},
		{
			name:"save tx_stake_delegate",
			args: struct{ tx store.Docs }{
				tx: doc_tx_stake_delegate,},

		},
		{
			name:"save tx_stake_unBond",
			args: struct{ tx store.Docs }{
				tx: doc_tx_stake_unbond,},

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveTx(tt.args.tx)
		})
	}
}

func Test_saveOrUpdateAccount(t *testing.T) {
	type args struct {
		tx store.Docs
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveOrUpdateAccount(tt.args.tx)
		})
	}
}

func Test_updateAccountBalance(t *testing.T) {
	type args struct {
		tx store.Docs
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateAccountBalance(tt.args.tx)
		})
	}
}
