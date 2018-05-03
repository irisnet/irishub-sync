package document

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNmStakeTxDeclareCandidacy = "tx_stake"
)

// Description
type Description struct {
	Moniker  string `bson:"moniker"`
	Identity string `bson:"identity"`
	Website  string `bson:"website"`
	Details  string `bson:"details"`
}

type StakeTxDeclareCandidacy struct {
	StakeTx `bson:"stake_tx"`
	Description `bson:"description"`
}

func (s StakeTxDeclareCandidacy) Name() string  {
	return CollectionNmStakeTxDeclareCandidacy
}

func (s StakeTxDeclareCandidacy) PkKvPair() map[string]interface{}  {
	return bson.M{"stake_tx.tx_hash": s.TxHash}
}

func (s StakeTxDeclareCandidacy) Index() []mgo.Index {
	return []mgo.Index{
		{
			Key:        []string{"description.moniker"},
			Unique:     false,
			DropDups:   false,
			Background: true,
		},
	}
}