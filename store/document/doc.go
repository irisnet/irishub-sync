package document

import (
	"github.com/irisnet/irishub-sync/store"
)

func init() {
	store.RegisterDocs(new(Account))
	store.RegisterDocs(new(Candidate))
	store.RegisterDocs(new(Delegator))
	store.RegisterDocs(new(Block))
	store.RegisterDocs(new(CommonTx))
	store.RegisterDocs(new(SyncTaskBak))
	store.RegisterDocs(new(ValidatorUpTime))
	store.RegisterDocs(new(TxGas))
	store.RegisterDocs(new(TxMsg))
	store.RegisterDocs(new(SyncTask))
	store.RegisterDocs(new(SyncConf))
}
