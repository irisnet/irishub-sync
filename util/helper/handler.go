package helper

import (
	"github.com/irisnet/irishub-sync/msg/nft"
	"github.com/irisnet/irishub-sync/msg/iservice"
	"github.com/irisnet/irishub-sync/msg/oracle"
	"github.com/irisnet/irishub-sync/msg/evidence"
	"github.com/irisnet/irishub-sync/msg/crisis"
	"github.com/irisnet/irishub-sync/msg/record"
	"github.com/irisnet/irishub-sync/msg/gov"
	"github.com/irisnet/irishub-sync/msg/htlc"
	"github.com/irisnet/irishub-sync/msg/coinswap"
	"github.com/irisnet/irishub-sync/msg/guardian"
	"github.com/irisnet/irishub-sync/msg/staking"
	"github.com/irisnet/irishub-sync/msg/distribution"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sync/msg/random"
	"github.com/irisnet/irishub-sync/msg/bank"
	"github.com/irisnet/irishub-sync/store/document"
)

func HandleMsg(msgData sdk.Msg, tx *document.CommonTx) (*document.CommonTx, bool) {
	if NftTx, ok := nft.HandleTxMsg(msgData, tx); ok {
		return NftTx, ok
	}
	if IserviceTx, ok := iservice.HandleTxMsg(msgData, tx); ok {
		return IserviceTx, ok
	}
	if OracleTx, ok := oracle.HandleTxMsg(msgData, tx); ok {
		return OracleTx, ok
	}
	if EvidenceTx, ok := evidence.HandleTxMsg(msgData, tx); ok {
		return EvidenceTx, ok
	}
	if CrisisTx, ok := crisis.HandleTxMsg(msgData, tx); ok {
		return CrisisTx, ok
	}
	if RecordTx, ok := record.HandleTxMsg(msgData, tx); ok {
		return RecordTx, ok
	}
	if GovTx, ok := gov.HandleTxMsg(msgData, tx); ok {
		return GovTx, ok
	}
	if HtlcTx, ok := htlc.HandleTxMsg(msgData, tx); ok {
		return HtlcTx, ok
	}
	if CoinswapTx, ok := coinswap.HandleTxMsg(msgData, tx); ok {
		return CoinswapTx, ok
	}
	if GuardianTx, ok := guardian.HandleTxMsg(msgData, tx); ok {
		return GuardianTx, ok
	}
	if StakeTx, ok := staking.HandleTxMsg(msgData, tx); ok {
		return StakeTx, ok
	}
	if DistriTx, ok := distribution.HandleTxMsg(msgData, tx); ok {
		return DistriTx, ok
	}
	if BankTx, ok := bank.HandleTxMsg(msgData, tx); ok {
		return BankTx, ok
	}
	if RandomTx, ok := random.HandleTxMsg(msgData, tx); ok {
		return RandomTx, ok
	}
	return tx, false
}
