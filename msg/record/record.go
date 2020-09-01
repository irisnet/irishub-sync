package record



import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/store/document"
)

func HandleTxMsg(msg sdk.Msg, tx *document.CommonTx) (*document.CommonTx, bool) {

	switch  msg.Type() {
	case new(types.MsgCreateRecord).Type():
		docMsg := DocMsgCreateRecord{}
		return docMsg.HandleTxMsg(msg, tx), true
	}
	return tx, false
}

