package server

import (
	"os"
	"strconv"
	"strings"

	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/util/constant"
)

var (
	BlockChainMonitorUrl = []string{"tcp://192.168.150.7:30657"}
	ChainId              = "rainbow-dev"

	WorkerNumCreateTask  = 2
	WorkerNumExecuteTask = 60

	InitConnectionNum        = 50              // fast init num of tendermint client pool
	MaxConnectionNum         = 100             // max size of tendermint client pool
	CronWatchBlock           = "*/1 * * * * *" // every 1 seconds
	CronCalculateUpTime      = "0 */1 * * * *" // every minute
	CronCalculateTxGas       = "0 */5 * * * *" // every five minute
	SyncProposalStatus       = "0 */1 * * * *" // every minute
	CronSaveValidatorHistory = "@daily"        // every day
	CronUpdateDelegator      = "0/5 * * * * *" // every ten minute

	// deprecated
	SyncMaxGoroutine = 60 // max go routine in server
	// deprecated
	SyncBlockNumFastSync = 8000 // sync block num each goroutine

	ConsulAddr    = "192.168.150.7:8500"
	SyncWithDLock = false
	Network       = "testnet"
)

// get value of env var
func init() {
	nodeUrl, found := os.LookupEnv(constant.EnvNameSerNetworkFullNode)
	if found {
		BlockChainMonitorUrl = strings.Split(nodeUrl, ",")
	}

	logger.Info("Env Value", logger.Any(constant.EnvNameSerNetworkFullNode, BlockChainMonitorUrl))

	chainId, found := os.LookupEnv(constant.EnvNameSerNetworkChainId)
	if found {
		ChainId = chainId
	}
	logger.Info("Env Value", logger.String(constant.EnvNameSerNetworkChainId, ChainId))

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
	logger.Info("Env Value", logger.String(constant.EnvNameCronSaveValidatorHistory, cronSaveValidatorHistory))

	workerNumCreateTask, found := os.LookupEnv(constant.EnvNameWorkerNumCreateTask)
	if found {
		var err error
		WorkerNumCreateTask, err = strconv.Atoi(workerNumCreateTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(constant.EnvNameWorkerNumCreateTask, workerNumCreateTask))
		}
	}
	logger.Info("Env Value", logger.Int(constant.EnvNameWorkerNumCreateTask, WorkerNumCreateTask))

	workerNumExecuteTask, found := os.LookupEnv(constant.EnvNameWorkerNumExecuteTask)
	if found {
		var err error
		WorkerNumExecuteTask, err = strconv.Atoi(workerNumExecuteTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(constant.EnvNameWorkerNumExecuteTask, workerNumCreateTask))
		}
	}
	logger.Info("Env Value", logger.Int(constant.EnvNameWorkerNumExecuteTask, WorkerNumExecuteTask))

	network, found := os.LookupEnv(constant.EnvNameNetwork)
	if found {
		Network = network
	}
	logger.Info("Env Value", logger.String(constant.EnvNameNetwork, network))
}
