package handler

import (
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/helper"
	"reflect"
	"sync"
)

// get tx type
func GetTxType(docTx store.Docs) string {
	if !reflect.ValueOf(docTx).FieldByName("Type").IsValid() {
		logger.Error.Printf("type which is field name of stake docTx is missed, docTx is %+v\n",
			helper.ToJson(docTx))
		return ""
	}
	txType := reflect.ValueOf(docTx).FieldByName("Type").String()

	return txType
}

func Handle(docTx store.Docs, mutex sync.Mutex, funChains []func(tx store.Docs, mutex sync.Mutex)) {
	for _, fun := range funChains {
		fun(docTx, mutex)
	}
}
