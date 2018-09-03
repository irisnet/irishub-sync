// interface for a document

package store

type Docs interface {
	// collection name
	Name() string
	// primary key pair(used to find a unique record)
	PkKvPair() map[string]interface{}
}

type Coin struct {
	Denom  string  `json:"denom"`
	Amount float64 `json:"amount"`
}

type Coins []Coin

type Fee struct {
	Amount Coins
	Gas    int64
}

type ActualFee struct {
	Denom  string  `json:"denom"`
	Amount float64 `json:"amount"`
}
