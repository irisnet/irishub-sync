package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/store/document"
)

func HandleTxMsg(msg sdk.Msg, tx *document.CommonTx) (*document.CommonTx, bool) {

	switch  msg.Type() {
	case new(types.MsgStartFeed).Type():
		docMsg := DocMsgStartFeed{}
		return docMsg.HandleTxMsg(msg, tx), true
	case new(types.MsgPauseFeed).Type():
		docMsg := DocMsgPauseFeed{}
		return docMsg.HandleTxMsg(msg, tx), true
	case new(types.MsgEditFeed).Type():
		docMsg := DocMsgEditFeed{}
		return docMsg.HandleTxMsg(msg, tx), true
	case new(types.MsgCreateFeed).Type():
		docMsg := DocMsgCreateFeed{}
		return docMsg.HandleTxMsg(msg, tx), true
	}
	return tx, false
}
