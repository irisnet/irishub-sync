// package for define constants

package constant

const (
	TxTypeBank  = "coin"
	TxTypeStakeCreate = "declareCandidacy"
	TxTypeStakeEdit = "editCandidacy"
	TxTypeStakeDelegate = "delegate"
	TxTypeStakeUnbond = "unbond"

	TxStatusSuccess = "success"
	
	EnvNameDbHost = "DB_HOST"
	EnvNameDbPort = "DB_PORT"
	EnvNameDbUser = "DB_USER"
	EnvNameDbPassWd = "DB_PASSWD"
	EnvNameDbDataBase = "DB_DATABASE"
	
	EnvNameSerNetworkNodeUrl = "SER_BC_NODE_URL"
	EnvNameSerNetworkChainId = "SER_BC_CHAIN_ID"
	EnvNameSerNetworkToken = "SER_BC_TOKEN"
	EnvNameSerMaxGoRoutine = "SER_MAX_GOROUTINE"
	EnvNameSerSyncBlockNum = "SER_SYNC_BLOCK_NUM"
)
