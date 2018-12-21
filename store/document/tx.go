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

	Tx_Field_Time                 = "time"
	Tx_Field_Height               = "height"
	Tx_Field_Hash                 = "tx_hash"
	Tx_Field_From                 = "from"
	Tx_Field_To                   = "to"
	Tx_Field_Amount               = "amount"
	Tx_Field_Type                 = "type"
	Tx_Field_Fee                  = "fee"
	Tx_Field_Memo                 = "memo"
	Tx_Field_Status               = "status"
	Tx_Field_Code                 = "code"
	Tx_Field_Log                  = "log"
	Tx_Field_GasUsed              = "gas_used"
	Tx_Field_GasPrice             = "gas_price"
	Tx_Field_ActualFee            = "actual_fee"
	Tx_Field_ProposalId           = "proposal_id"
	Tx_Field_Tags                 = "tags"
	Tx_Field_StakeCreateValidator = "stake_create_validator"
	Tx_Field_StakeEditValidator   = "stake_edit_validator"
)

type CommonTx struct {
	Time       time.Time         `bson:"time"`
	Height     int64             `bson:"height"`
	TxHash     string            `bson:"tx_hash"`
	From       string            `bson:"from"`
	To         string            `bson:"to"`
	Amount     store.Coins       `bson:"amount"`
	Type       string            `bson:"type"`
	Fee        store.Fee         `bson:"fee"`
	Memo       string            `bson:"memo"`
	Status     string            `bson:"status"`
	Code       uint32            `bson:"code"`
	Log        string            `bson:"log"`
	GasUsed    int64             `bson:"gas_used"`
	GasPrice   float64           `bson:"gas_price"`
	ActualFee  store.ActualFee   `bson:"actual_fee"`
	ProposalId uint64            `bson:"proposal_id"`
	Tags       map[string]string `bson:"tags"`

	StakeCreateValidator StakeCreateValidator `bson:"stake_create_validator"`
	StakeEditValidator   StakeEditValidator   `bson:"stake_edit_validator"`
	Msg                  store.Msg            `bson:"-"`
}

// Description
type ValDescription struct {
	Moniker  string `bson:"moniker"`
	Identity string `bson:"identity"`
	Website  string `bson:"website"`
	Details  string `bson:"details"`
}

type StakeCreateValidator struct {
	PubKey      string         `bson:"pub_key"`
	Description ValDescription `bson:"description"`
}

type StakeEditValidator struct {
	Description ValDescription `bson:"description"`
}

func (d CommonTx) Name() string {
	return CollectionNmCommonTx
}

func (d CommonTx) PkKvPair() map[string]interface{} {
	return bson.M{Tx_Field_Hash: d.TxHash}
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
