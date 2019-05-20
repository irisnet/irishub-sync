// package for define constants

package constant

const (
	TxTypeTransfer                    = "Transfer"
	TxTypeBurn                        = "Burn"
	TxTypeStakeCreateValidator        = "CreateValidator"
	TxTypeStakeEditValidator          = "EditValidator"
	TxTypeStakeDelegate               = "Delegate"
	TxTypeStakeBeginUnbonding         = "BeginUnbonding"
	TxTypeBeginRedelegate             = "BeginRedelegate"
	TxTypeUnjail                      = "Unjail"
	TxTypeSetWithdrawAddress          = "SetWithdrawAddress"
	TxTypeWithdrawDelegatorReward     = "WithdrawDelegatorReward"
	TxTypeWithdrawDelegatorRewardsAll = "WithdrawDelegatorRewardsAll"
	TxTypeWithdrawValidatorRewardsAll = "WithdrawValidatorRewardsAll"
	TxTypeSubmitProposal              = "SubmitProposal"
	TxTypeDeposit                     = "Deposit"
	TxTypeVote                        = "Vote"

	EnvNameDbAddr     = "DB_ADDR"
	EnvNameDbUser     = "DB_USER"
	EnvNameDbPassWd   = "DB_PASSWD"
	EnvNameDbDataBase = "DB_DATABASE"

	EnvNameSerNetworkFullNode   = "SER_BC_FULL_NODE"
	EnvNameSerNetworkChainId    = "SER_BC_CHAIN_ID"
	EnvNameWorkerNumCreateTask  = "WORKER_NUM_CREATE_TASK"
	EnvNameWorkerNumExecuteTask = "WORKER_NUM_EXECUTE_TASK"

	EnvNameNetwork = "NETWORK"

	EnvLogFileName    = "LOG_FILE_NAME"
	EnvLogFileMaxSize = "LOG_FILE_MAX_SIZE"
	EnvLogFileMaxAge  = "LOG_FILE_MAX_AGE"
	EnvLogCompress    = "LOG_COMPRESS"
	EnableAtomicLevel = "ENABLE_ATOMIC_LEVEL"

	// define store name
	StoreNameStake      = "stake"
	StoreDefaultEndPath = "key"

	StatusDepositPeriod = "DepositPeriod"
	StatusVotingPeriod  = "VotingPeriod"
	StatusPassed        = "Passed"
	StatusRejected      = "Rejected"

	NetworkMainnet = "mainnet"

	// define coin unit
	IrisAttoUnit = "iris-atto"
)
