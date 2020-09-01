package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmCommonTx = "tx_common"
	TxStatusSuccess      = "success"
	TxStatusFail         = "fail"
	Unknow_Status        = "unknown"

	Tx_Field_Hash   = "tx_hash"
	Tx_Field_Type   = "type"
	Tx_Field_Status = "status"
	Tx_Field_Height = "height"
)

type CommonTx struct {
	Time       time.Time       `bson:"time"`
	Height     int64           `bson:"height"`
	TxHash     string          `bson:"tx_hash"`
	From       string          `bson:"from"`
	To         string          `bson:"to"`
	Amount     store.Coins     `bson:"amount"`
	Type       string          `bson:"type"`
	Fee        store.Fee       `bson:"fee"`
	Memo       string          `bson:"memo"`
	Status     string          `bson:"status"`
	Code       uint32          `bson:"code"`
	Log        string          `bson:"log"`
	GasUsed    int64           `bson:"gas_used"`
	GasWanted  int64           `bson:"gas_wanted"`
	GasPrice   float64         `bson:"gas_price"`
	ActualFee  store.ActualFee `bson:"actual_fee"`
	ProposalId uint64          `bson:"proposal_id"`
	Events     []Event         `bson:"events"`

	Signers  []Signer `bson:"signers"`
	Signer   string   `bson:"signer"`
	Addrs    []string `bson:"addrs"`
	Types    []string `bson:"types"`
	TimeUnix int64    `bson:"time_unix"`

	Msgs []DocTxMsg `bson:"msgs"`
}

type DocTxMsg struct {
	Type string `bson:"type"`
	Msg  Msg    `bson:"msg"`
}

type Event struct {
	Type       string      `bson:"type" json:"type"`
	Attributes []Attribute `bson:"attributes" json:"attributes"`
}

type Attribute struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value" json:"value"`
}
type Msg interface {
	Type() string
	BuildMsg(msg interface{})
}

// Description
type ValDescription struct {
	Moniker  string `bson:"moniker"`
	Identity string `bson:"identity"`
	Website  string `bson:"website"`
	Details  string `bson:"details"`
}

type CommissionMsg struct {
	Rate          string `bson:"rate"`            // the commission rate charged to delegators
	MaxRate       string `bson:"max_rate"`        // maximum commission rate which validator can ever charge
	MaxChangeRate string `bson:"max_change_rate"` // maximum daily increase of the validator commission
}

type Signer struct {
	AddrHex    string `bson:"addr_hex"`
	AddrBech32 string `bson:"addr_bech32"`
}

func (d CommonTx) Name() string {
	return CollectionNmCommonTx
}

func (d CommonTx) PkKvPair() map[string]interface{} {
	return bson.M{Tx_Field_Hash: d.TxHash}
}
func (d CommonTx) EnsureIndexs() []mgo.Index {
	var indexes []mgo.Index
	indexes = append(indexes, mgo.Index{
		Key:        []string{"-tx_hash"},
		Unique:     true,
		Background: true,
	})

	return indexes
}

func (d CommonTx) Query(query, fields bson.M, sort []string, skip, limit int) (
	results []CommonTx, err error) {
	exop := func(c *mgo.Collection) error {
		return c.Find(query).Sort(sort...).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	return results, store.ExecCollection(d.Name(), exop)
}

func (d CommonTx) CalculateTxGasAndGasPrice(txType string, limit int) (
	[]CommonTx, error) {
	query := bson.M{
		Tx_Field_Type:   txType,
		Tx_Field_Status: TxStatusSuccess,
	}
	fields := bson.M{}
	sort := []string{"-height"}
	skip := 0

	return d.Query(query, fields, sort, skip, limit)
}

func (d CommonTx) GetUnknownOrEmptyTypeTxs(skip, limit int) (res []CommonTx, err error) {
	q := bson.M{"$or": []bson.M{
		{Tx_Field_Status: Unknow_Status},
		{Tx_Field_Type: ""},
	}}
	sorts := []string{"-height"}
	selector := bson.M{
		Tx_Field_Hash:   1,
		Tx_Field_Height: 1,
	}

	fn := func(c *mgo.Collection) error {
		return c.Find(q).Select(selector).Sort(sorts...).Skip(skip).Limit(limit).All(&res)
	}

	err = store.ExecCollection(CollectionNmCommonTx, fn)
	if err != nil {
		return nil, err
	}
	return
}
