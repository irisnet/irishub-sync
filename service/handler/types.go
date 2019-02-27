package handler

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"sync"
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
