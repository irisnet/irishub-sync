package handler

import (
	"strings"

	"github.com/irisnet/irishub-sync/store/document"
	itypes "github.com/irisnet/irishub-sync/types"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

func handleTokenFlow(blockWithtags document.Block, tx document.CommonTx) []txn.Op {
	result := []txn.Op{}
	for _, v := range blockWithtags.Result.EndBlock.Tags {
		if strings.ToUpper(v.Key) == tx.TxHash {
			tokenFlow, ok := parseTagsAndTx(tx, v.Value)
			if ok {
				tokenFlow.BlockHeight = blockWithtags.Height
				tokenFlow.BlockHash = blockWithtags.Hash
				txOp := txn.Op{
					C:      document.CollectionNmCommonTokenFlow,
					Id:     bson.NewObjectId(),
					Insert: tokenFlow,
				}
				result = append(result, txOp)
			}
		}
	}
	return result
}

func parseTagsAndTx(tx document.CommonTx, tagStr string) (document.CommonTokenFlow, bool) {
	var result document.CommonTokenFlow
	fromToValue := strings.Split(tagStr, "::")
	if len(fromToValue) != 6 {
		return result, false
	}
	result.TxHash = tx.TxHash
	result.From = fromToValue[0]
	result.To = fromToValue[1]
	result.FlowType = fromToValue[3]
	result.Amount = itypes.ParseCoin(fromToValue[2])
	result.TxInitiator = tx.From
	result.Timestamp = fromToValue[5]
	result.Fee = tx.Fee
	result.Status = tx.Status
	result.TxType = tx.Type
	return result, true
}
