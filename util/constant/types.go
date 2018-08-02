// package for define constants

package constant

const (
	TxTypeTransfer               = "transfer"
	TxTypeStakeCreateValidator   = "createValidator"
	TxTypeStakeEditValidator     = "editValidator"
	TxTypeStakeDelegate          = "delegate"
	TxTypeStakeBeginUnbonding    = "beginUnbonding"
	TxTypeStakeCompleteUnbonding = "completeUnbonding"

	TxStatusSuccess = "success"
	TxStatusFail    = "fail"

	EnvNameDbHost     = "DB_HOST"
	EnvNameDbPort     = "DB_PORT"
	EnvNameDbUser     = "DB_USER"
	EnvNameDbPassWd   = "DB_PASSWD"
	EnvNameDbDataBase = "DB_DATABASE"

	EnvNameSerNetworkNodeUrl = "SER_BC_NODE_URL"
	EnvNameSerNetworkChainId = "SER_BC_CHAIN_ID"
	EnvNameSerNetworkToken   = "SER_BC_TOKEN"
	EnvNameSerMaxGoRoutine   = "SER_MAX_GOROUTINE"
	EnvNameSerSyncBlockNum   = "SER_SYNC_BLOCK_NUM"

	// define store name
	StoreNameStake      = "stake"
	StoreDefaultEndPath = "key"

	// define sync type
	SyncTypeFastSync = "fastSync"
	SyncTypeWatch    = "watch"
)
