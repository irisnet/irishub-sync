package handler

import (
	"github.com/irisnet/irishub-sync/store/document"
	"sync"
)

// get tx type
func GetTxType(docTx document.CommonTx) string {
	if docTx.TxHash == "" {
		return ""
	}
	txType := docTx.Type

	return txType
}

func Handle(docTx document.CommonTx, mutex sync.Mutex, funChains []func(tx document.CommonTx, mutex sync.Mutex)) {
	for _, fun := range funChains {
		if docTx.TxHash != "" {
			fun(docTx, mutex)
		}
	}
}
