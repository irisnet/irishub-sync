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
	return docTx.Type
}

type Action = func(tx document.CommonTx, mutex sync.Mutex)

func Handle(docTx document.CommonTx, mutex sync.Mutex, actions []Action) {
	for _, action := range actions {
		if docTx.TxHash != "" {
			action(docTx, mutex)
		}
	}
}
