// package for define constants

package constant

const (
	TxTypeTransfer               = "Transfer"
	TxTypeStakeCreateValidator   = "CreateValidator"
	TxTypeStakeEditValidator     = "EditValidator"
	TxTypeStakeDelegate          = "Delegate"
	TxTypeStakeBeginUnbonding    = "BeginUnbonding"
	TxTypeStakeCompleteUnbonding = "CompleteUnbonding"

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

	// define interval block num and tx num
	IntervalBlockNumCalculateValidatorUpTime = int64(100)
	IntervalTxNumCalculateTxGas              = 100
)
