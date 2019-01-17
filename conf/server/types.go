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

	Bech32 = Bech32AddrPrefix{
		PrefixAccAddr:  "faa",
		PrefixAccPub:   "fap",
		PrefixValAddr:  "fva",
		PrefixValPub:   "fvp",
		PrefixConsAddr: "fca",
		PrefixConsPub:  "fcp",
	}
)

type Bech32AddrPrefix struct {
	PrefixAccAddr  string
	PrefixAccPub   string
	PrefixValAddr  string
	PrefixValPub   string
	PrefixConsAddr string
	PrefixConsPub  string
}

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

	loadBe32Prefix()
}

func loadBe32Prefix() {
	prefixAccAddr, found := os.LookupEnv(constant.EnvNamePrefixAccAddr)
	if found {
		Bech32.PrefixAccAddr = prefixAccAddr
	}
	logger.Info("Env Value", logger.String(constant.EnvNamePrefixAccAddr, Bech32.PrefixAccAddr))

	prefixAccPub, found := os.LookupEnv(constant.EnvNamePrefixAccPub)
	if found {
		Bech32.PrefixAccPub = prefixAccPub
	}
	logger.Info("Env Value", logger.String(constant.EnvNamePrefixAccPub, Bech32.PrefixAccPub))

	prefixValAddr, found := os.LookupEnv(constant.EnvNamePrefixValAddr)
	if found {
		Bech32.PrefixValAddr = prefixValAddr
	}
	logger.Info("Env Value", logger.String(constant.EnvNamePrefixValAddr, Bech32.PrefixValAddr))

	prefixValPub, found := os.LookupEnv(constant.EnvNamePrefixValPub)
	if found {
		Bech32.PrefixValPub = prefixValPub
	}
	logger.Info("Env Value", logger.String(constant.EnvNamePrefixValPub, Bech32.PrefixValPub))

	prefixConsAddr, found := os.LookupEnv(constant.EnvNamePrefixConsAddr)
	if found {
		Bech32.PrefixConsAddr = prefixConsAddr
	}
	logger.Info("Env Value", logger.String(constant.EnvNamePrefixConsAddr, Bech32.PrefixConsAddr))

	prefixConsPub, found := os.LookupEnv(constant.EnvNamePrefixConsPub)
	if found {
		Bech32.PrefixConsPub = prefixConsPub
	}
	logger.Info("Env Value", logger.String(constant.EnvNamePrefixConsPub, Bech32.PrefixConsPub))
}
