package handler

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"sync"
	"github.com/irisnet/irishub-sync/store"
)

// get tx type
func GetTxType(docTx document.CommonTx) string {
	if docTx.TxHash == "" {
		return ""
	}
	return docTx.Type
}

type Action = func(tx document.CommonTx, mutex sync.Mutex)

func Handle(docTx document.CommonTx, mutex sync.Mutex, actions []Action) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Parse Tx failed", logger.Int64("height", docTx.Height),
				logger.String("txHash", docTx.TxHash), logger.Any("err", err))
		}
	}()

	for _, action := range actions {
		if docTx.TxHash != "" {
			action(docTx, mutex)
		}
	}
}

var (
	SyncTaskModel document.SyncTask
	BlockModel    document.Block
	TxModel       document.CommonTx
	Account       document.Account
	Proposal      document.Proposal
	SyncConf      document.SyncConf

	Collections = []store.Docs{
		SyncTaskModel,
		BlockModel,
		TxModel,
		Account,
		Proposal,
		SyncConf,
	}
)

func EnsureDocsIndexes() {
	if len(Collections) > 0 {
		for _, v := range Collections {
			if indexs := v.EnsureIndexs(); len(indexs) > 0 {
				store.EnsureIndexes(v.Name(), indexs)
			}
		}
	}
}
