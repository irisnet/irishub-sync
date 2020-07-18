package msg

import (
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util/constant"
)

type (
	DocTxMsgIssueToken struct {
		Symbol        string `bson:"symbol"`
		Name          string `bson:"name"`
		Scale         uint32 `bson:"scale"`
		MinUnit       string `bson:"min_unit"`
		InitialSupply uint64 `bson:"initial_supply"`
		MaxSupply     uint64 `bson:"max_supply"`
		Mintable      bool   `bson:"mintable"`
		Owner         string `bson:"owner"`
	}

	DocTxMsgEditToken struct {
		Owner     string `bson:"owner"` //  owner of token
		MaxSupply uint64 `bson:"max_supply"`
		Mintable  bool   `bson:"mintable"` //  mintable of token
		Symbol    string `bson:"symbol"`
		Name      string `bson:"name"`
	}

	DocTxMsgMintToken struct {
		Symbol string `bson:"symbol"`
		Owner  string `bson:"owner"`  // the current owner address of the token
		To     string `bson:"to"`     // address of mint token to
		Amount uint64 `bson:"amount"` // amount of mint token
	}

	DocTxMsgTransferTokenOwner struct {
		SrcOwner string `bson:"src_owner"` // the current owner address of the token
		DstOwner string `bson:"dst_owner"` // the new owner
		Symbol   string `bson:"symbol"`
	}
)

func (m *DocTxMsgIssueToken) Type() string {
	return constant.TxMsgTypeAssetIssueToken
}

func (m *DocTxMsgIssueToken) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgIssueToken)

	m.Symbol = msg.Symbol
	m.Name = msg.Name
	m.Scale = msg.Scale
	m.MinUnit = msg.MinUnit
	m.InitialSupply = msg.InitialSupply
	m.MaxSupply = msg.MaxSupply
	m.Mintable = msg.Mintable
	m.Owner = msg.Owner.String()
}

func (m *DocTxMsgEditToken) Type() string {
	return constant.TxMsgTypeAssetEditToken
}

func (m *DocTxMsgEditToken) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgEditToken)

	m.Symbol = msg.Symbol
	m.Owner = msg.Owner.String()
	m.MaxSupply = msg.MaxSupply
	switch msg.Mintable {
	case constant.TrueStr:
		m.Mintable = true
		break
	case constant.FalseStr:
		m.Mintable = false
		break
	}
	m.Name = msg.Name
}

func (m *DocTxMsgMintToken) Type() string {
	return constant.TxMsgTypeAssetMintToken
}

func (m *DocTxMsgMintToken) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgMintToken)

	m.Symbol = msg.Symbol
	m.Owner = msg.Owner.String()
	m.To = msg.To.String()
	m.Amount = msg.Amount
}

func (m *DocTxMsgTransferTokenOwner) Type() string {
	return constant.TxMsgTypeAssetTransferTokenOwner
}

func (m *DocTxMsgTransferTokenOwner) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgTransferTokenOwner)

	m.SrcOwner = msg.SrcOwner.String()
	m.DstOwner = msg.DstOwner.String()
	m.Symbol = msg.Symbol
}
