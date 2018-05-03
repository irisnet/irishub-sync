package document

import (
	"github.com/irisnet/iris-sync-server/model/store"
)

func init() {
	store.RegisterDocs(new(Account))
	store.RegisterDocs(new(CoinTx))
	store.RegisterDocs(new(StakeTx))
	store.RegisterDocs(new(StakeTxDeclareCandidacy))
	store.RegisterDocs(new(Candidate))
	store.RegisterDocs(new(Delegator))
	store.RegisterDocs(new(Block))
	store.RegisterDocs(new(SyncTask))
}
