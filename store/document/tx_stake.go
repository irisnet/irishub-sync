package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmStakeTx                 = "tx_stake"
	CollectionNmStakeTxDeclareCandidacy = CollectionNmStakeTx
	CollectionNmStakeTxEditCandidacy    = CollectionNmStakeTx
)

// ===============================
// struct of delegate and unbond
// ===============================
type StakeTx struct {
	TxHash        string     `bson:"tx_hash"`
	Time          time.Time  `bson:"time"`
	Height        int64      `bson:"height"`
	DelegatorAddr string     `bson:"from"`
	ValidatorAddr string     `bson:"to"`
	PubKey        string     `bson:"pub_key"`
	Type          string     `bson:"type"`
	Amount        store.Coin `bson:"amount"`
	Fee           store.Fee  `bson:"fee"`
	Status        string     `bson:"status"`
}

// Description
type Description struct {
	Moniker  string `bson:"moniker"`
	Identity string `bson:"identity"`
	Website  string `bson:"website"`
	Details  string `bson:"details"`
}

func (c StakeTx) Name() string {
	return CollectionNmStakeTx
}

func (c StakeTx) PkKvPair() map[string]interface{} {
	return bson.M{"tx_hash": c.TxHash}
}

// ============================
// struct of createValidator
// ============================

type StakeTxDeclareCandidacy struct {
	StakeTx     `bson:"stake_tx"`
	Description `bson:"description"`
}

func (s StakeTxDeclareCandidacy) Name() string {
	return CollectionNmStakeTxDeclareCandidacy
}

func (s StakeTxDeclareCandidacy) PkKvPair() map[string]interface{} {
	return bson.M{"stake_tx.tx_hash": s.TxHash}
}

// ============================
// struct of editValidator
// ============================

type StakeTxEditCandidacy struct {
	StakeTx     `bson:"stake_tx"`
	Description `bson:"description"`
}

func (s StakeTxEditCandidacy) Name() string {
	return CollectionNmStakeTxEditCandidacy
}

func (s StakeTxEditCandidacy) PkKvPair() map[string]interface{} {
	return bson.M{"stake_tx.tx_hash": s.TxHash}
}
