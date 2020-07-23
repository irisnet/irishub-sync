package nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/store/document"
)

func HandleTxMsg(msg sdk.Msg, tx *document.CommonTx) (*document.CommonTx, bool) {

	switch  msg.Type() {
	case new(types.MsgMintNFT).Type():
		docMsg := DocMsgNFTMint{}
		return docMsg.HandleTxMsg(msg, tx), true
	case new(types.MsgEditNFT).Type():
		docMsg := DocMsgNFTEdit{}
		return docMsg.HandleTxMsg(msg, tx), true
	case new(types.MsgTransferNFT).Type():
		docMsg := DocMsgNFTTransfer{}
		return docMsg.HandleTxMsg(msg, tx), true
	case new(types.MsgBurnNFT).Type():
		docMsg := DocMsgNFTBurn{}
		return docMsg.HandleTxMsg(msg, tx), true
	case new(types.MsgIssueDenom).Type():
		docMsg := DocMsgIssueDenom{}
		return docMsg.HandleTxMsg(msg, tx), true
	}
	return tx, false
}
