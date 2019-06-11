package server

import (
	"os"
	"strconv"
	"strings"

	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/util/constant"
)

var (
	BlockChainMonitorUrl = []string{"tcp://35.201.147.145:30657"}

	WorkerNumCreateTask  = 2
	WorkerNumExecuteTask = 60

	InitConnectionNum  = 50              // fast init num of tendermint client pool
	MaxConnectionNum   = 100             // max size of tendermint client pool
	SyncProposalStatus = "0 */1 * * * *" // every minute

	Network = "testnet"
)

// get value of env var
func init() {
	nodeUrl, found := os.LookupEnv(constant.EnvNameSerNetworkFullNode)
	if found {
		BlockChainMonitorUrl = strings.Split(nodeUrl, ",")
	}

	logger.Info("Env Value", logger.Any(constant.EnvNameSerNetworkFullNode, BlockChainMonitorUrl))

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
	logger.Info("Env Value", logger.String(constant.EnvNameNetwork, Network))
}
