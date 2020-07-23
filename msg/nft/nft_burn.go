package nft

import (
	. "github.com/irisnet/irishub-sync/util/constant"
	"strings"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
	"github.com/irisnet/irishub-sync/store"
)

type (
	DocMsgNFTBurn struct {
		Sender string `bson:"sender"`
		ID     string `bson:"id"`
		Denom  string `bson:"denom"`
	}
)

func (m *DocMsgNFTBurn) Type() string {
	return TxTypeNFTBurn
}

func (m *DocMsgNFTBurn) BuildMsg(v interface{}) {
	var msg types.MsgBurnNFT
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.Sender = msg.Sender.String()
	m.ID = strings.ToLower(msg.ID)
	m.Denom = strings.ToLower(msg.Denom)
}

func (m *DocMsgNFTBurn) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Type = m.Type()
	tx.From = m.Sender
	tx.To = ""
	tx.Amount = []store.Coin{}
	return tx
}
