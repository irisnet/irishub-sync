package msg

import (
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store"
	itypes "github.com/irisnet/irishub-sync/types"
	"encoding/hex"
)

type DocTxMsgCreateHTLC struct {
	Sender               string      `bson:"sender"`                  // the initiator address
	To                   string      `bson:"to"`                      // the destination address
	ReceiverOnOtherChain string      `bson:"receiver_on_other_chain"` // the claim receiving address on the other chain
	Amount               store.Coins `bson:"amount"`                  // the amount to be transferred
	HashLock             string      `bson:"hash_lock"`               // the hash lock generated from secret (and timestamp if provided)
	Timestamp            uint64      `bson:"timestamp"`               // if provided, used to generate the hash lock together with secret
	TimeLock             uint64      `bson:"time_lock"`               // the time span after which the HTLC will expire
}

func (doctx *DocTxMsgCreateHTLC) Type() string {
	return constant.TxTypeCreateHTLC
}

func (doctx *DocTxMsgCreateHTLC) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgCreateHTLC)
	doctx.Sender = msg.Sender.String()
	doctx.To = msg.To.String()
	doctx.Amount = itypes.ParseCoins(msg.Amount.String())
	doctx.Timestamp = msg.Timestamp
	doctx.HashLock = hex.EncodeToString(msg.HashLock)
	doctx.TimeLock = msg.TimeLock
	doctx.ReceiverOnOtherChain = msg.ReceiverOnOtherChain
}

type DocTxMsgClaimHTLC struct {
	Sender   string `bson:"sender"`    // the initiator address
	HashLock string `bson:"hash_lock"` // the hash lock identifying the HTLC to be claimed
	Secret   string `bson:"secret"`    // the secret with which to claim
}

func (doctx *DocTxMsgClaimHTLC) Type() string {
	return constant.TxTypeClaimHTLC
}

func (doctx *DocTxMsgClaimHTLC) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgClaimHTLC)
	doctx.Sender = msg.Sender.String()
	doctx.Secret = hex.EncodeToString(msg.Secret)
	doctx.HashLock = hex.EncodeToString(msg.HashLock)
}

type DocTxMsgRefundHTLC struct {
	Sender   string `bson:"sender"`    // the initiator address
	HashLock string `bson:"hash_lock"` // the hash lock identifying the HTLC to be refunded
}

func (doctx *DocTxMsgRefundHTLC) Type() string {
	return constant.TxTypeRefundHTLC
}

func (doctx *DocTxMsgRefundHTLC) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgRefundHTLC)
	doctx.Sender = msg.Sender.String()
	doctx.HashLock = hex.EncodeToString(msg.HashLock)
}
