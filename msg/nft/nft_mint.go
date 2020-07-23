package nft

import (
	. "github.com/irisnet/irishub-sync/util/constant"
	"strings"
	"github.com/irisnet/irishub-sync/store/document"
	"encoding/json"
	"github.com/irisnet/irishub-sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/store"
)

type DocMsgNFTMint struct {
	Sender    string `bson:"sender"`
	Recipient string `bson:"recipient"`
	Denom     string `bson:"denom"`
	ID        string `bson:"id"`
	TokenURI  string `bson:"token_uri"`
	TokenData string `bson:"token_data"`
}

func (m *DocMsgNFTMint) Type() string {
	return TxTypeNFTMint
}

func (m *DocMsgNFTMint) BuildMsg(v interface{}) {
	var msg types.MsgMintNFT
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.Sender = msg.Sender.String()
	m.Recipient = msg.Recipient.String()
	m.ID = strings.ToLower(msg.ID)
	m.Denom = strings.ToLower(msg.Denom)
	m.TokenURI = msg.TokenURI
	m.TokenData = msg.TokenData
}

func (m *DocMsgNFTMint) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.From = m.Sender
	tx.To = m.Recipient
	tx.Amount = []store.Coin{}
	tx.Type = m.Type()
	return tx
}
