// interface for a document

package store

const (
	CollectionNameTxn = "mgo_txn"
)

type Docs interface {
	// collection name
	Name() string
	// primary key pair(used to find a unique record)
	PkKvPair() map[string]interface{}
}

type Coin struct {
	Denom  string  `json:"denom" bson:"denom"`
	Amount float64 `json:"amount" bson:"amount"`
}

type Coins []Coin

type Fee struct {
	Amount Coins `json:"amount"`
	Gas    int64 `json:"gas"`
}

type ActualFee struct {
	Denom  string  `json:"denom"`
	Amount float64 `json:"amount"`
}

type Msg interface {
	Type() string
	String() string
}
