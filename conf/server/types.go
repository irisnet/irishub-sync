package server

import (
	"os"
	"strconv"
	"strings"

	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/util/constant"
)

var (
	BlockChainMonitorUrl = []string{"tcp://127.0.0.1:26657"}
	ChainId              = "rainbow-dev"
	Token                = "iris"

	InitConnectionNum        = 50              // fast init num of tendermint client pool
	MaxConnectionNum         = 100             // max size of tendermint client pool
	CronWatchBlock           = "*/1 * * * * *" // every 1 seconds
	CronCalculateUpTime      = "0 */1 * * * *" // every minute
	CronCalculateTxGas       = "0 */5 * * * *" // every five minute
	SyncProposalStatus       = "0 */1 * * * *" // every minute
	CronSaveValidatorHistory = "@daily"        // every day

	SyncMaxGoroutine     = 60   // max go routine in server
	SyncBlockNumFastSync = 8000 // sync block num each goroutine
	ConsulAddr           = "192.168.150.7:8500"
	SyncWithDLock        = false
)

// get value of env var
func init() {
	nodeUrl, found := os.LookupEnv(constant.EnvNameSerNetworkNodeUrl)
	if found {
		BlockChainMonitorUrl = strings.Split(nodeUrl, ",")
	}

	logger.Info("Env Value", logger.Any(constant.EnvNameSerNetworkNodeUrl, BlockChainMonitorUrl))

	chainId, found := os.LookupEnv(constant.EnvNameSerNetworkChainId)
	if found {
		ChainId = chainId
	}
	logger.Info("Env Value", logger.String(constant.EnvNameSerNetworkChainId, ChainId))

	token, found := os.LookupEnv(constant.EnvNameSerNetworkToken)
	if found {
		Token = token
	}
	logger.Info("Env Value", logger.String(constant.EnvNameSerNetworkToken, Token))

	maxGoroutine, found := os.LookupEnv(constant.EnvNameSerMaxGoRoutine)
	if found {
		var err error
		SyncMaxGoroutine, err = strconv.Atoi(maxGoroutine)
		if err != nil {
			logger.Fatal("Env Value", logger.String(constant.EnvNameSerMaxGoRoutine, maxGoroutine))
		}
	}
	logger.Info("Env Value", logger.Int(constant.EnvNameSerMaxGoRoutine, SyncMaxGoroutine))

	syncBlockNum, found := os.LookupEnv(constant.EnvNameSerSyncBlockNum)
	if found {
		var err error
		SyncBlockNumFastSync, err = strconv.Atoi(syncBlockNum)
		if err != nil {
			logger.Fatal("Env Value", logger.String(constant.EnvNameSerSyncBlockNum, syncBlockNum))
		}
	}
	logger.Info("Env Value", logger.Int(constant.EnvNameSerSyncBlockNum, SyncBlockNumFastSync))

	consulAddr, found := os.LookupEnv(constant.EnvNameConsulAddr)
	if found {
		ConsulAddr = consulAddr
	}
	logger.Info("Env Value", logger.String(constant.EnvNameConsulAddr, ConsulAddr))

	withDLock, found := os.LookupEnv(constant.EnvNameSyncWithDLock)
	if found {
		flag, err := strconv.ParseBool(withDLock)
		if err != nil {
			logger.Fatal("Env Value", logger.String(constant.EnvNameSyncWithDLock, withDLock))
		}
		SyncWithDLock = flag
	}
	logger.Info("Env Value", logger.Bool(constant.EnvNameSyncWithDLock, SyncWithDLock))

	cronSaveValidatorHistory, found := os.LookupEnv(constant.EnvNameCronSaveValidatorHistory)
	if found {
		CronSaveValidatorHistory = cronSaveValidatorHistory
	}
	logger.Info("Env Value", logger.String(constant.EnvNameCronSaveValidatorHistory, ConsulAddr))

}
