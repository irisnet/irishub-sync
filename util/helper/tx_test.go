package helper

import (
	"testing"

	"github.com/irisnet/iris-sync-server/module/logger"

	_ "github.com/cosmos/cosmos-sdk/modules/auth"
	_ "github.com/cosmos/cosmos-sdk/modules/base"

	"github.com/tendermint/tendermint/types"
)

func buildTxByte(blockHeight int64) types.Tx {

	InitClientPool()
	client := GetClient()
	// release client
	defer client.Release()

	block, err := client.Client.Block(&blockHeight)

	if err != nil {
		logger.Error.Panic(err)
	}

	if block.BlockMeta.Header.NumTxs > 0 {
		txs := block.Block.Data.Txs
		return txs[0]
	}

	return nil
}

func TestParseTx(t *testing.T) {

	txCoinByte := buildTxByte(12453)
	txStakeDeclareCandidacyByte := buildTxByte(19073)
	txStakeDelegateByte := buildTxByte(13725)
	txStakeUnBondByte := buildTxByte(14260)


	type args struct {
		txByte types.Tx
	}
	tests := []struct {
		name  string
		args  args
	}{
		{
			name: "tx_coin_send",
			args: struct{ txByte types.Tx }{
				txByte: txCoinByte,
					},
		},
		{
			name: "tx_stake_declareCandidacy",
			args: struct{ txByte types.Tx }{
				txByte: txStakeDeclareCandidacyByte,
					},
		},
		{
			name: "tx_stake_delegate",
			args: struct{ txByte types.Tx }{
				txByte: txStakeDelegateByte,
			},
		},

		{
			name: "tx_stake_unbond",
			args: struct{ txByte types.Tx }{
				txByte: txStakeUnBondByte,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txType, txContent := ParseTx(tt.args.txByte)
			logger.Info.Printf("%s: tx type is %s, and struct is %+v\n", tt.name, txType, txContent)
		})
	}
}
