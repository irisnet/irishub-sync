package record

import (
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/store/document"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/store"
)

type (
	DocMsgCreateRecord struct {
		Contents []Content `bson:"contents"`
		Creator  string    `bson:"creator"`
	}

	Content struct {
		Digest     string `bson:"digest"`
		DigestAlgo string `bson:"digest_algo"`
		URI        string `bson:"uri"`
		Meta       string `bson:"meta"`
	}
)

func (d *DocMsgCreateRecord) Type() string {
	return constant.TxTypeCreateRecord
}

func (d *DocMsgCreateRecord) BuildMsg(msg interface{}) {
	m := msg.(*types.MsgCreateRecord)

	var docContents []Content
	if len(m.Contents) > 0 {
		for _, v := range m.Contents {
			docContents = append(docContents, Content{
				Digest:     v.Digest,
				DigestAlgo: v.DigestAlgo,
				URI:        v.URI,
				Meta:       v.Meta,
			})
		}
	}

	d.Contents = docContents
	d.Creator = m.Creator.String()
}

func (m *DocMsgCreateRecord) HandleTxMsg(msgData sdk.Msg, tx *document.CommonTx) *document.CommonTx {

	m.BuildMsg(msgData)
	tx.Msgs = append(tx.Msgs, document.DocTxMsg{
		Type: m.Type(),
		Msg:  m,
	})
	tx.Addrs = append(tx.Addrs, m.Creator)
	tx.Types = append(tx.Types, m.Type())
	if len(tx.Msgs) > 1 {
		return tx
	}
	tx.Type = m.Type()
	tx.From = m.Creator
	tx.To = ""
	tx.Amount = []store.Coin{}

	return tx
}

//func (d *DocMsgRecordC) HandleTxMsg(msg types.MsgCreateRecord) MsgDocInfo {
//	var (
//		from, to, signer string
//		coins            []store.Coin
//		docTxMsg         document.DocTxMsg
//		complexMsg       bool
//		signers          []string
//		addrs            []string
//	)
//
//	from = msg.Creator.String()
//	to = ""
//
//	d.BuildMsg(msg)
//	docTxMsg = document.DocTxMsg{
//		Type: d.Type(),
//		Msg:  d,
//	}
//	complexMsg = false
//	signer, signers = store.BuildDocSigners(msg.GetSigners())
//	addrs = append(addrs, signers...)
//	addrs = append(addrs, d.Creator)
//
//	return MsgDocInfo{
//		From:       from,
//		To:         to,
//		Coins:      coins,
//		Signer:     signer,
//		DocTxMsg:   docTxMsg,
//		ComplexMsg: complexMsg,
//		Signers:    signers,
//		Addrs:      addrs,
//	}
//}
