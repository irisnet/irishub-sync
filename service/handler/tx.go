package handler

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"sync"
)

// save Tx document into collection
func SaveTx(docTx document.CommonTx, mutex sync.Mutex) {
	var (
		methodName = "SaveTx"
	)
	logger.Debug("Start", logger.String("method", methodName))

	// save common docTx document
	saveCommonTx := func(commonTx document.CommonTx) {
		//save tx
		err := store.Save(commonTx)
		if err != nil {
			logger.Error("Save commonTx failed", logger.Any("Tx", commonTx), logger.String("err", err.Error()))
		}
		//save tx_msg
		msg := commonTx.Msg
		if msg != nil {
			txMsg := document.TxMsg{
				Hash:    docTx.TxHash,
				Type:    msg.Type(),
				Content: msg.String(),
			}
			store.Save(txMsg)
		}
		handleProposal(commonTx)
	}

	saveCommonTx(docTx)
	logger.Debug("End", logger.String("method", methodName))
}
