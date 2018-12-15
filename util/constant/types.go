// package for define constants

package constant

const (
	TxTypeTransfer                    = "Transfer"
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
	TxTypeStakeCompleteUnbonding      = "CompleteUnbonding"
	TxTypeSubmitProposal              = "SubmitProposal"
	TxTypeDeposit                     = "Deposit"
	TxTypeVote                        = "Vote"


	EnvNameDbAddr     = "DB_ADDR"
	EnvNameDbUser     = "DB_USER"
	EnvNameDbPassWd   = "DB_PASSWD"
	EnvNameDbDataBase = "DB_DATABASE"

	EnvNameSerNetworkNodeUrl        = "SER_BC_NODE_URL"
	EnvNameSerNetworkChainId        = "SER_BC_CHAIN_ID"
	EnvNameSerNetworkToken          = "SER_BC_TOKEN"
	EnvNameSerMaxGoRoutine          = "SER_MAX_GOROUTINE"
	EnvNameSerSyncBlockNum          = "SER_SYNC_BLOCK_NUM"
	EnvNameConsulAddr               = "CONSUL_ADDR"
	EnvNameSyncWithDLock            = "SYNC_WITH_DLOCK"
	EnvNameCronSaveValidatorHistory = "CRON_SAVE_VALIDATOR_HISTORY"

	//bech32Prefix
	EnvNamePrefixAccAddr  = "PrefixAccAddr"
	EnvNamePrefixAccPub   = "PrefixAccPub"
	EnvNamePrefixValAddr  = "PrefixValAddr"
	EnvNamePrefixValPub   = "PrefixValPub"
	EnvNamePrefixConsAddr = "PrefixConsAddr"
	EnvNamePrefixConsPub  = "PrefixConsPub"

	EnvLogFileName    = "LOG_FILE_NAME"
	EnvLogFileMaxSize = "LOG_FILE_MAX_SIZE"
	EnvLogFileMaxAge  = "LOG_FILE_MAX_AGE"
	EnvLogCompress    = "LOG_COMPRESS"
	EnableAtomicLevel = "ENABLE_ATOMIC_LEVEL"

	// define store name
	StoreNameStake      = "stake"
	StoreDefaultEndPath = "key"

	// define sync type
	SyncTypeFastSync = "fastSync"
	SyncTypeWatch    = "watch"

	// define interval block num and tx num
	IntervalBlockNumCalculateValidatorUpTime = int64(100)
	IntervalTxNumCalculateTxGas              = 100

	StatusDepositPeriod = "DepositPeriod"
	StatusVotingPeriod  = "VotingPeriod"
	StatusPassed        = "Passed"
	StatusRejected      = "Rejected"
)
