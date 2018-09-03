package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/irisnet/irishub-sync/module/logger"
)

const CollectionNmTxGas = "tx_gas"

type TxGas struct {
	TxType   string   `bson:"tx_type"`
	GasUsed  GasUsed  `bson:"gas_used"`
	GasPrice GasPrice `bson:"gas_price"`
}

type GasUsed struct {
	MinGasUsed float64 `bson:"min_gas_used"`
	MaxGasUsed float64 `bson:"max_gas_used"`
	AvgGasUsed float64 `bson:"avg_gas_used"`
}

type GasPrice struct {
	Denom       string  `bson:"denom"`
	MinGasPrice float64 `bson:"min_gas_price"`
	MaxGasPrice float64 `bson:"max_gas_price"`
	AvgGasPrice float64 `bson:"avg_gas_price"`
}

func (d TxGas) Name() string {
	return CollectionNmTxGas
}

func (d TxGas) PkKvPair() map[string]interface{} {
	return bson.M{"tx_type": d.TxType}
}

func (d TxGas) RemoveAll() error {
	query := bson.M{}
	remove := func(c *mgo.Collection) error {
		changeInfo, err := c.RemoveAll(query)
		logger.Info.Printf("remove all tx gas data, %+v", changeInfo)
		return err
	}
	return store.ExecCollection(d.Name(), remove)
}

func (d TxGas) SaveAll(txGases []TxGas) error {
	var docs []interface{}

	if len(txGases) == 0 {
		return nil
	}

	for _, v := range txGases {
		docs = append(docs, v)
	}

	err := store.SaveAll(d.Name(), docs)

	return err
}
