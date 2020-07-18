package nft

import (
	. "github.com/irisnet/irishub-sync/util/constant"
	"strings"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
)

type (
	DocMsgNFTEdit struct {
		Sender    string `bson:"sender"`
		ID        string `bson:"id"`
		Denom     string `bson:"denom"`
		TokenURI  string `bson:"token_uri"`
		TokenData string `bson:"token_data"`
	}
)

func (m *DocMsgNFTEdit) Type() string {
	return MsgTypeNFTEdit
}

func (m *DocMsgNFTEdit) BuildMsg(v interface{}) {
	var msg types.MsgEditNFT
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.Sender = msg.Sender.String()
	m.ID = strings.ToLower(msg.ID)
	m.Denom = strings.ToLower(msg.Denom)
	m.TokenURI = msg.TokenURI
	m.TokenData = msg.TokenData
}

func (m *DocMsgNFTEdit) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Type = m.Type()
	return tx
}
