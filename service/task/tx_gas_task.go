package task

import (
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
)

func calculateTxGasAndGasPrice() {
	var (
		methodName    = "CalculateTxGasAndGasPrice"
		intervalTxNum = constant.IntervalTxNumCalculateTxGas
		txModel       document.CommonTx
		txGasModel    document.TxGas
		txGases       []document.TxGas
	)
	logger.Info("Start", logger.String("method", methodName))

	txTypes := []string{
		constant.TxTypeTransfer,
		constant.TxTypeStakeCreateValidator,
		constant.TxTypeStakeEditValidator,
		constant.TxTypeStakeDelegate,
		constant.TxTypeStakeBeginUnbonding,
		constant.TxTypeBeginRedelegate,
		constant.TxTypeSetWithdrawAddress,
		constant.TxTypeWithdrawDelegatorReward,
		constant.TxTypeWithdrawDelegatorRewardsAll,
		constant.TxTypeWithdrawValidatorRewardsAll,
	}

	for _, v := range txTypes {
		txs, err := txModel.CalculateTxGasAndGasPrice(v, intervalTxNum)
		if err != nil {
			logger.Error("Can't calculate gas and gasPrice", logger.String("err", err.Error()))
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
		logger.Error("Remove all data fail", logger.String("err", err.Error()))
		return
	}

	// save all data
	err2 := txGasModel.SaveAll(txGases)
	if err2 != nil {
		logger.Error("Save latest data fail", logger.String("err", err2.Error()))
		return
	}

	logger.Info("End", logger.String("method", methodName))
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
			maxGasUsed = v.GasUsed
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

func MakeCalculateTxGasAndGasPriceTask() Task {
	return NewLockTaskFromEnv(conf.CronCalculateTxGas, "calculate_tx_gas_and_gas_price_lock", func() {
		logger.Debug("========================task's trigger [CalculateTxGasAndGasPrice] begin===================")
		calculateTxGasAndGasPrice()
		logger.Debug("========================task's trigger [CalculateTxGasAndGasPrice] end===================")
	})
}
