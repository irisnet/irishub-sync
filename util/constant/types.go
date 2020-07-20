// package for define constants

package constant

const (
	TxTypeTransfer                       = "Transfer"
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
	TxTypeMsgFundCommunityPool           = "FundCommunityPool "
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

	MsgTypeNFTMint     = "NFTMint"
	MsgTypeNFTEdit     = "NFTEdit"
	MsgTypeNFTTransfer = "NFTTransfer"
	MsgTypeNFTBurn     = "NFTBurn"
	MsgTypeIssueDenom  = "IssueDenom"

	MsgTypeDefineService         = "DefineService"         // type for MsgDefineService
	MsgTypeBindService           = "BindService"           // type for MsgBindService
	MsgTypeUpdateServiceBinding  = "UpdateServiceBinding"  // type for MsgUpdateServiceBinding
	MsgTypeSetWithdrawAddress    = "SetWithdrawAddress"    // type for MsgSetWithdrawAddress
	MsgTypeDisableServiceBinding = "DisableServiceBinding" // type for MsgDisableServiceBinding
	MsgTypeEnableServiceBinding  = "EnableServiceBinding"  // type for MsgEnableServiceBinding
	MsgTypeRefundServiceDeposit  = "RefundServiceDeposit"  // type for MsgRefundServiceDeposit
	MsgTypeCallService           = "CallService"           // type for MsgCallService
	MsgTypeRespondService        = "RespondService"        // type for MsgRespondService
	MsgTypePauseRequestContext   = "PauseRequestContext"   // type for MsgPauseRequestContext
	MsgTypeStartRequestContext   = "StartRequestContex"    // type for MsgStartRequestContext
	MsgTypeKillRequestContext    = "KillRequestContex"     // type for MsgKillRequestContext
	MsgTypeUpdateRequestContext  = "UpdateRequestContext"  // type for MsgUpdateRequestContext
	MsgTypeWithdrawEarnedFees    = "WithdrawEarnedFees"    // type for MsgWithdrawEarnedFees

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

	TxMsgTypeSubmitProposal                = "SubmitProposal"
	TxMsgTypeSubmitSoftwareUpgradeProposal = "SubmitSoftwareUpgradeProposal"
	TxMsgTypeSubmitTaxUsageProposal        = "SubmitTaxUsageProposal"
	TxMsgTypeSubmitTokenAdditionProposal   = "SubmitTokenAdditionProposal"
	TxMsgTypeAssetIssueToken               = "IssueToken"
	TxMsgTypeAssetEditToken                = "EditToken"
	TxMsgTypeAssetMintToken                = "MintToken"
	TxMsgTypeAssetTransferTokenOwner       = "TransferTokenOwner"
	TxMsgTypeAssetCreateGateway            = "CreateGateway"
	TxMsgTypeAssetEditGateway              = "EditGateway"
	TxMsgTypeAssetTransferGatewayOwner     = "TransferGatewayOwner"

	TxTypeSubmitEvidence  = "SubmitEvidence"
	TxTypeVerifyInvariant = "VerifyInvariant"

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
