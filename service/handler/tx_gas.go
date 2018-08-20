package handler

import (
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
)

func CalculateTxGasAndGasPrice() {
	var (
		methodName    = "CalculateTxGasAndGasPrice"
		intervalTxNum = constant.IntervalTxNumCalculateTxGas
		txModel       document.CommonTx
		txGasModel    document.TxGas
		txGases       []document.TxGas
	)
	logger.Info.Printf("%v: Start\n", methodName)

	txTypes := []string{
		constant.TxTypeTransfer,
		constant.TxTypeStakeCreateValidator,
		constant.TxTypeStakeEditValidator,
		constant.TxTypeStakeDelegate,
		constant.TxTypeStakeBeginUnbonding,
		constant.TxTypeStakeCompleteUnbonding,
	}

	for _, v := range txTypes {
		txs, err := txModel.CalculateTxGasAndGasPrice(v, intervalTxNum)
		if err != nil {
			logger.Error.Printf("%v: Can't calculate gas and gasPrice, err is %v\n",
				methodName, err)
			continue
		}
		if len(txs) > 0 {
			txGas := buildTxGas(txs)
			txGases = append(txGases, txGas)
		}
	}

	// remove all data
	err := txGasModel.RemoveAll()
	if err != nil {
		logger.Error.Printf("%v: Remove all data fail, err is %v\n",
			methodName, err)
		return
	}

	// save all data
	err2 := txGasModel.SaveAll(txGases)
	if err2 != nil {
		logger.Error.Printf("%v: Save latest data fail, err is %v\n",
			methodName, err2)
		return
	}

	logger.Info.Printf("%v: End\n", methodName)
}

func buildTxGas(txs []document.CommonTx) document.TxGas {
	var (
		txGas         document.TxGas
		gasPriceDenom string
	)

	// get type of tx
	txType := txs[0].Type
	// get denom of gasPrice
	// all gasPrice denom are the same at version v0.23.0-iris1,
	// so can set first fee denom as gasPrice denom。
	// these code should be refactored at next version
	if len(txs[0].Fee.Amount) > 0 {
		gasPriceDenom = txs[0].Fee.Amount[0].Denom
	}
	minGasUsed := txs[0].GasUsed
	maxGasUsed := txs[0].GasUsed
	totalGasUsed := float64(0)
	minGasPrice := txs[0].GasPrice
	maxGasPrice := txs[0].GasPrice
	totalGasPrice := float64(0)
	for _, v := range txs {
		if v.GasUsed < minGasUsed {
			minGasUsed = v.GasUsed
		}
		if v.GasUsed > maxGasUsed {
			minGasUsed = v.GasUsed
		}
		totalGasUsed += float64(v.GasUsed)

		if v.GasPrice < minGasPrice {
			minGasPrice = v.GasPrice
		}
		if v.GasPrice > maxGasPrice {
			maxGasPrice = v.GasPrice
		}
		totalGasPrice += v.GasPrice
	}

	txGas = document.TxGas{
		TxType: txType,
		GasUsed: document.GasUsed{
			MinGasUsed: float64(minGasUsed),
			MaxGasUsed: float64(maxGasUsed),
			AvgGasUsed: totalGasUsed / float64(len(txs)),
		},
		GasPrice: document.GasPrice{
			Denom:       gasPriceDenom,
			MinGasPrice: minGasPrice,
			MaxGasPrice: maxGasPrice,
			AvgGasPrice: totalGasPrice / float64(len(txs)),
		},
	}

	return txGas
}
