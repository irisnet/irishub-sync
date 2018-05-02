package document

import (
	"gopkg.in/mgo.v2"
)

const (
	CollectionNmStakeTxDeclareCandidacy = "tx_stake"
)

type StakeTxDeclareCandidacy struct {
	StakeTx
	Description
}

// Stake交易
type Description struct {
	Moniker  string `bson:"moniker"`
	Identity string `bson:"identity"`
	Website  string `bson:"website"`
	Details  string `bson:"details"`
}

func Name() string  {
	return CollectionNmStakeTxDeclareCandidacy
}

func PkKvPair() map[string]interface{}  {
	return nil
}

func Index() mgo.Index {
	return mgo.Index{}
}