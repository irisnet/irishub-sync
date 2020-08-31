// package for define constants

package constant

const (
	TxTypeTransfer                       = "Transfer"
	TxTypeMultiSend                      = "MultiSend"
	TxTypeBurn                           = "Burn"
	TxTypeSetMemoRegexp                  = "SetMemoRegexp"
	TxTypeStakeCreateValidator           = "CreateValidator"
	TxTypeStakeEditValidator             = "EditValidator"
	TxTypeStakeDelegate                  = "Delegate"
	TxTypeStakeBeginUnbonding            = "BeginUnbonding"
	TxTypeBeginRedelegate                = "BeginRedelegate"
	TxTypeUnjail                         = "Unjail"
	TxTypeSetWithdrawAddress             = "SetWithdrawAddress"
	TxTypeWithdrawDelegatorReward        = "WithdrawDelegatorReward"
	TxTypeMsgFundCommunityPool           = "FundCommunityPool"
	TxTypeMsgWithdrawValidatorCommission = "WithdrawValidatorCommission"
	TxTypeSubmitProposal                 = "SubmitProposal"
	TxTypeDeposit                        = "Deposit"
	TxTypeVote                           = "Vote"
	TxTypeRequestRand                    = "RequestRand"
	TxTypeAssetIssueToken                = "IssueToken"
	TxTypeAssetEditToken                 = "EditToken"
	TxTypeAssetMintToken                 = "MintToken"
	TxTypeAssetTransferTokenOwner        = "TransferTokenOwner"
	TxTypeAssetCreateGateway             = "CreateGateway"
	TxTypeAssetEditGateway               = "EditGateway"
	TxTypeAssetTransferGatewayOwner      = "TransferGatewayOwner"

	TxTypeNFTMint     = "NFTMint"
	TxTypeNFTEdit     = "NFTEdit"
	TxTypeNFTTransfer = "NFTTransfer"
	TxTypeNFTBurn     = "NFTBurn"
	TxTypeIssueDenom  = "IssueDenom"

	TxTypeDefineService             = "DefineService"              // type for MsgDefineService
	TxTypeBindService               = "BindService"                // type for MsgBindService
	TxTypeUpdateServiceBinding      = "UpdateServiceBinding"       // type for MsgUpdateServiceBinding
	TxTypeServiceSetWithdrawAddress = "service/SetWithdrawAddress" // type for SetWithdrawFeesAddress
	TxTypeDisableServiceBinding     = "DisableServiceBinding"      // type for MsgDisableServiceBinding
	TxTypeEnableServiceBinding      = "EnableServiceBinding"       // type for MsgEnableServiceBinding
	TxTypeRefundServiceDeposit      = "RefundServiceDeposit"       // type for MsgRefundServiceDeposit
	TxTypeCallService               = "CallService"                // type for MsgCallService
	TxTypeRespondService            = "RespondService"             // type for MsgRespondService
	TxTypePauseRequestContext       = "PauseRequestContext"        // type for MsgPauseRequestContext
	TxTypeStartRequestContext       = "StartRequestContext"        // type for MsgStartRequestContext
	TxTypeKillRequestContext        = "KillRequestContext"         // type for MsgKillRequestContext
	TxTypeUpdateRequestContext      = "UpdateRequestContext"       // type for MsgUpdateRequestContext
	TxTypeWithdrawEarnedFees        = "WithdrawEarnedFees"         // type for MsgWithdrawEarnedFees

	TxTypeAddProfiler    = "AddProfiler"
	TxTypeAddTrustee     = "AddTrustee"
	TxTypeDeleteTrustee  = "DeleteTrustee"
	TxTypeDeleteProfiler = "DeleteProfiler"

	TxTypeCreateFeed = "CreateFeed"
	TxTypeEditFeed   = "EditFeed"
	TxTypePauseFeed  = "PauseFeed"
	TxTypeStartFeed  = "StartFeed"

	TxTypeCreateHTLC = "CreateHTLC"
	TxTypeClaimHTLC  = "ClaimHTLC"
	TxTypeRefundHTLC = "RefundHTLC"

	TxTypeAddLiquidity    = "AddLiquidity"
	TxTypeRemoveLiquidity = "RemoveLiquidity"
	TxTypeSwapOrder       = "SwapOrder"

	TxTypeSubmitEvidence  = "SubmitEvidence"
	TxTypeVerifyInvariant = "VerifyInvariant"

	MsgTypeRecordCreate = "CreateRecord"

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

	TrueStr  = "true"
	FalseStr = "false"
)
