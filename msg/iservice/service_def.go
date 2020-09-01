package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/types"
	. "github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store/document"
	"encoding/json"
	"github.com/irisnet/irishub-sync/store"
)

type (
	DocMsgDefineService struct {
		Name              string   `bson:"name" yaml:"name"`
		Description       string   `bson:"description" yaml:"description"`
		Tags              []string `bson:"tags" yaml:"tags"`
		Author            string   `bson:"author" yaml:"author"`
		AuthorDescription string   `bson:"author_description" yaml:"author_description"`
		Schemas           string   `bson:"schemas"`
	}
)

func (m *DocMsgDefineService) Type() string {
	return TxTypeDefineService
}

func (m *DocMsgDefineService) BuildMsg(v interface{}) {
	var msg types.MsgDefineService
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &msg)

	m.Name = msg.Name
	m.Description = msg.Description
	m.Tags = msg.Tags
	m.Author = msg.Author.String()
	m.AuthorDescription = msg.AuthorDescription
	m.Schemas = msg.Schemas
}

func (m *DocMsgDefineService) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Addrs = append(tx.Addrs, m.Author)
	tx.Types = append(tx.Types, m.Type())
	if len(tx.Msgs) > 1 {
		return tx
	}
	tx.Type = m.Type()
	tx.From = m.Author
	tx.To = ""
	tx.Amount = []store.Coin{}
	return tx
}
