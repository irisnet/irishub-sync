package handler

import (
	"testing"

	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/helper"
	"os"
)

const (
	BankHeight                   = 17694
	StakeCreateHeight            = 28848
	StakeEditHeight              = 28581
	StakeDelegateHeight          = 79026
	StakeBeginUnbondingHeight    = 79063
	StakeCompleteUnbondingHeight = 79177
)

func TestMain(m *testing.M) {
	// setup
	store.Start()

	code := m.Run()

	// shutdown
	os.Exit(code)
}

func buildDocData(blockHeight int64) document.CommonTx {
	client := helper.GetClient()
	// release client
	defer client.Release()

	block, err := client.Client.Block(&blockHeight)

	if err != nil {
		logger.Panic(err.Error())
	}

	if block.BlockMeta.Header.NumTxs > 0 {
		txs := block.Block.Data.Txs
		txByte := txs[0]
		docTx,_ := helper.ParseTx(txByte, block.Block)

		return docTx

	}
	return document.CommonTx{}
}
