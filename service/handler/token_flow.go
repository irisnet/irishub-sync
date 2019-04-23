package handler

import (
	"strings"
	"time"

	"github.com/irisnet/irishub-sync/store/document"
	itypes "github.com/irisnet/irishub-sync/types"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

const (
	TokenFlowTagsLen = 6
)

func handleTokenFlow(blockWithTags document.Block, tx document.CommonTx, result *[]txn.Op) {
	for _, v := range blockWithTags.Result.EndBlock.Tags {
		if strings.ToUpper(v.Key) == tx.TxHash {
			tokenFlow, ok := parseTagsAndTx(tx, v.Value)
			if ok {
				tokenFlow.BlockHeight = blockWithTags.Height
				tokenFlow.BlockHash = blockWithTags.Hash
				txOp := txn.Op{
					C:      document.CollectionNmTokenFlow,
					Id:     bson.NewObjectId(),
					Insert: tokenFlow,
				}
				*result = append(*result, txOp)
			}
		}
	}
}

func parseTagsAndTx(tx document.CommonTx, tagStr string) (document.TokenFlow, bool) {
	var result document.TokenFlow
	flowItem := strings.Split(tagStr, "::")
	if len(flowItem) != TokenFlowTagsLen {
		return result, false
	}
	result.TxHash = tx.TxHash
	result.From = flowItem[0]
	result.To = flowItem[1]
	result.FlowType = flowItem[3]
	result.Amount = itypes.ParseCoin(flowItem[2])
	result.TxInitiator = tx.From
	t, err := time.Parse("2006-01-02 15:04:05.999999999 +0000 UTC", flowItem[5])
	if err != nil {
		panic(err.Error())
	}
	result.Timestamp = t
	result.Fee = tx.Fee
	result.Status = tx.Status
	result.TxType = tx.Type
	return result, true
}
