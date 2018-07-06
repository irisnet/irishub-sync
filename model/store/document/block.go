package document

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNmBlock = "block"
)

type Block struct {
	Height      int64             `bson:"height"`
	Hash        string            `bson:"hash"`
	Time        time.Time         `bson:"time"`
	NumTxs      int64             `bson:"num_txs"`
	Meta        BlockMeta `bson:"meta"`
	Block BlockContent `bson:"block"`
	Validators []Validator `bson:"validators"`
}

type BlockMeta struct {
	BlockID BlockID `bson:"block_id"`
	Header Header `bson:"header"`
}

type BlockID struct {
	Hash        string  `bson:"hash"`
	PartsHeader PartSetHeader `bson:"parts"`
}

type PartSetHeader struct {
	Total int    `bson:"total"`
	Hash  string `bson:"hash"`
}

type Header struct {
	// basic block info
	ChainID string    `bson:"chain_id"`
	Height  int64     `bson:"height"`
	Time    time.Time `bson:"time"`
	NumTxs  int64     `bson:"num_txs"`

	// prev block info
	LastBlockID BlockID `bson:"last_block_id"`
	TotalTxs    int64   `bson:"total_txs"`

	// hashes of block data
	LastCommitHash string `bson:"last_commit_hash"` // commit from validators from the last block
	DataHash       string `bson:"data_hash"`        // transactions

	// hashes from the app output from the prev block
	ValidatorsHash  string `bson:"validators_hash"`   // validators for the current block
	ConsensusHash   string `bson:"consensus_hash"`    // consensus params for current block
	AppHash         string `bson:"app_hash"`          // state after txs from the previous block
	LastResultsHash string `bson:"last_results_hash"` // root hash of all results from the txs from the previous block

	// consensus info
	EvidenceHash string `bson:"evidence_hash"` // evidence included in the block
}

type BlockContent struct {
	LastCommit Commit      `bson:"last_commit"`
}

type Commit struct {
	// NOTE: The Precommits are in order of address to preserve the bonded ValidatorSet order.
	// Any peer with a block can gossip precommits by index with a peer without recalculating the
	// active ValidatorSet.
	BlockID    BlockID `bson:"block_id"`
	Precommits []Vote `bson:"precommits"`
}

// Represents a prevote, precommit, or commit vote from validators for consensus.
type Vote struct {
	ValidatorAddress string          `bson:"validator_address"`
	ValidatorIndex   int              `bson:"validator_index"`
	Height           int64            `bson:"height"`
	Round            int              `bson:"round"`
	Timestamp        time.Time        `bson:"timestamp"`
	Type             byte             `bson:"type"`
	BlockID          BlockID          `bson:"block_id"` // zero if vote is nil.
	Signature        Signature `bson:"signature"`
}

type Signature struct {
	Type string `bson:"type"`
	Value string `bson:"value"`
}

type Validator struct {
	Address     string       `bson:"address"`
	PubKey      string `bson:"pub_key"`
	VotingPower int64         `bson:"voting_power"`
	Accum int64 `bson:"accum"`
} 


func (d Block) Name() string {
	return CollectionNmBlock
}

func (d Block) PkKvPair() map[string]interface{} {
	return bson.M{"height": d.Height}
}