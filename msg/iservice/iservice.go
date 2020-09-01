package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/store"
)

func HandleTxMsg(msg sdk.Msg, tx *document.CommonTx) (*document.CommonTx, bool) {
	ok := true
	switch msg.Type() {
	case new(types.MsgDefineService).Type():
		docMsg := DocMsgDefineService{}
		return docMsg.HandleTxMsg(msg, tx), ok
	case new(types.MsgBindService).Type():
		docMsg := DocMsgBindService{}
		return docMsg.HandleTxMsg(msg, tx), ok
	case new(types.MsgCallService).Type():
		docMsg := DocMsgCallService{}
		return docMsg.HandleTxMsg(msg, tx), ok

	case new(types.MsgRespondService).Type():
		docMsg := DocMsgServiceResponse{}
		return docMsg.HandleTxMsg(msg, tx), ok

	case new(types.MsgUpdateServiceBinding).Type():
		docMsg := DocMsgUpdateServiceBinding{}
		return docMsg.HandleTxMsg(msg, tx), ok

	case new(types.MsgSetWithdrawAddress).Type():
		docMsg := DocMsgServiceSetWithdrawAddress{}
		return docMsg.HandleTxMsg(msg, tx)

	case new(types.MsgDisableServiceBinding).Type():
		docMsg := DocMsgDisableServiceBinding{}
		return docMsg.HandleTxMsg(msg, tx), ok

	case new(types.MsgEnableServiceBinding).Type():
		docMsg := DocMsgEnableServiceBinding{}
		return docMsg.HandleTxMsg(msg, tx), ok

	case new(types.MsgRefundServiceDeposit).Type():
		docMsg := DocMsgRefundServiceDeposit{}
		return docMsg.HandleTxMsg(msg, tx), ok

	case new(types.MsgUpdateRequestContext).Type():
		docMsg := DocMsgUpdateRequestContext{}
		return docMsg.HandleTxMsg(msg, tx), ok

	case new(types.MsgPauseRequestContext).Type():
		docMsg := DocMsgPauseRequestContext{}
		return docMsg.HandleTxMsg(msg, tx), ok

	case new(types.MsgStartRequestContext).Type():
		docMsg := DocMsgStartRequestContext{}
		return docMsg.HandleTxMsg(msg, tx), ok

	case new(types.MsgKillRequestContext).Type():
		docMsg := DocMsgKillRequestContext{}
		return docMsg.HandleTxMsg(msg, tx), ok

	case new(types.MsgWithdrawEarnedFees).Type():
		docMsg := DocMsgWithdrawEarnedFees{}
		return docMsg.HandleTxMsg(msg, tx), ok
	default:
		ok = false
	}
	return tx, ok
}

type Coin struct {
	Denom  string `json:"denom" bson:"denom"`
	Amount string `json:"amount" bson:"amount"`
}

type Coins []Coin

func (c Coins) Convert() (result store.Coins) {
	for _, val := range c {
		result = append(result, types.ParseCoin(val.Amount+val.Denom))
	}
	return
}

type Fee struct {
	Amount Coins `json:"amount"`
	Gas    int64 `json:"gas"`
}
