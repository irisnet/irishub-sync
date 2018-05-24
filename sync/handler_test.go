package sync

import (
	"testing"
	"github.com/irisnet/iris-sync-server/model/store"
	"github.com/irisnet/iris-sync-server/util/helper"
	"github.com/irisnet/iris-sync-server/module/logger"
	"github.com/irisnet/iris-sync-server/util/constant"
	"github.com/irisnet/iris-sync-server/model/store/document"
	"github.com/irisnet/iris-sync-server/module/stake"
	"sync"
)

func init()  {
	helper.InitClientPool()
	store.Init()
}

func buildDocData(blockHeight int64) store.Docs {

	client := helper.GetClient()
	// release client
	defer client.Release()

	block, err := client.Client.Block(&blockHeight)

	if err != nil {
		logger.Error.Panic(err)
	}

	if block.BlockMeta.Header.NumTxs > 0 {
		txs := block.Block.Data.Txs
		txByte := txs[0]
		txType, tx := helper.ParseTx(txByte)

		switch txType {
		case constant.TxTypeCoin:
			coinTx, _ := tx.(document.CoinTx)
			coinTx.Height = block.Block.Height
			coinTx.Time = block.Block.Time
			return coinTx
		case stake.TypeTxDeclareCandidacy:
			stakeTxDeclareCandidacy, _ := tx.(document.StakeTxDeclareCandidacy)
			stakeTxDeclareCandidacy.Height = block.Block.Height
			stakeTxDeclareCandidacy.Time = block.Block.Time
			return stakeTxDeclareCandidacy
		case stake.TypeTxEditCandidacy:
			break
		case stake.TypeTxDelegate, stake.TypeTxUnbond:
			stakeTx, _ := tx.(document.StakeTx)
			stakeTx.Height = block.Block.Height
			stakeTx.Time = block.Block.Time
			return stakeTx
		}

	}
	return nil
}

func Test_saveTx(t *testing.T) {

	docTxCoin := buildDocData(1756)
	docTxStakeDeclareCandidacy := buildDocData(1921)
	docTxStakeDeclareCandidacy2 := buildDocData(1927)
	docTxStakeDelegate := buildDocData(115135)
	docTxStakeUnBond := buildDocData(397201)

	type args struct {
		tx store.Docs
		mutex sync.Mutex
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name:"save tx_coin",
			args: args{
				tx:docTxCoin,
				mutex:sync.Mutex{},
			},
		},
		{
			name:"save tx_stake_declareCandidacy",
			args: args{
				tx:docTxStakeDeclareCandidacy,
				mutex:sync.Mutex{},
			},

		},
		{
			name:"save tx_stake_declareCandidacy2",
			args: args{
				tx:docTxStakeDeclareCandidacy2,
				mutex:sync.Mutex{},
			},
		},
		{
			name:"save tx_stake_delegate",
			args:args{
				tx:docTxStakeDelegate,
				mutex:sync.Mutex{},
			},

		},
		{
			name:"save tx_stake_unBond",
			args: args{
				tx: docTxStakeUnBond,
				mutex: sync.Mutex{},
			},

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveTx(tt.args.tx, tt.args.mutex)
		})
	}
}

func Test_saveOrUpdateAccount(t *testing.T) {
	
	docTxCoin := buildDocData(1756)
	docTxStakeDeclareCandidacy := buildDocData(1921)
	docTxStakeDeclareCandidacy2 := buildDocData(1927)
	docTxStakeDelegate := buildDocData(115135)
	docTxStakeUnBond := buildDocData(397201)

	type args struct {
		tx store.Docs
		mutex sync.Mutex
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name:"save tx_coin",
			args: args{
				tx: docTxCoin,
			},
		},
		{
			name:"save tx_stake_declareCandidacy",
			args: args{
				tx: docTxStakeDeclareCandidacy,
			},

		},
		{
			name:"save tx_stake_declareCandidacy2",
			args: args{
				tx: docTxStakeDeclareCandidacy2,
			},
			
		},
		{
			name:"save tx_stake_delegate",
			args: args{
				tx:docTxStakeDelegate,
			},

		},
		{
			name:"save tx_stake_unBond",
			args:args{
				tx:docTxStakeUnBond,
			},

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveOrUpdateAccount(tt.args.tx, tt.args.mutex)
		})
	}
}

func Test_updateAccountBalance(t *testing.T) {
	
	docTxCoin := buildDocData(1756)
	docTxStakeDeclareCandidacy := buildDocData(1921)
	docTxStakeDeclareCandidacy2 := buildDocData(1927)
	docTxStakeDelegate := buildDocData(115135)
	docTxStakeUnBond := buildDocData(397201)

	type args struct {
		tx store.Docs
		mutex sync.Mutex
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name:"tx_coin",
			args: args{
				tx:docTxCoin,
			},
		},
		{
			name:"tx_stake_declareCandidacy",
			args: args{
				tx:docTxStakeDeclareCandidacy,
			},

		},
		{
			name:"tx_stake_declareCandidacy2",
			args: args{
				tx:docTxStakeDeclareCandidacy2,
			},
			
		},
		{
			name:"tx_stake_delegate",
			args: args{
				tx:docTxStakeDelegate,
			},

		},
		{
			name:"tx_stake_unBond",
			args: args{
				tx:docTxStakeUnBond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateAccountBalance(tt.args.tx, tt.args.mutex)
		})
	}
}
