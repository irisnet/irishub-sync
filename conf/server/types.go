package server

import (
	"os"
	"strconv"

	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/util/constant"
)

var (
	BlockChainMonitorUrl = "tcp://127.0.0.1:26657"
	ChainId              = "rainbow-dev"
	Token                = "iris"

	InitConnectionNum   = 50              // fast init num of tendermint client pool
	MaxConnectionNum    = 100             // max size of tendermint client pool
	CronWatchBlock      = "*/1 * * * * *" // every 10 seconds
	CronCalculateUpTime = "0 */1 * * * *" // every minute
	CronCalculateTxGas  = "0 */5 * * * *" // every five minute
	SyncProposalStatus  = "0 */1 * * * *" // every minute

	SyncMaxGoroutine     = 60   // max go routine in server
	SyncBlockNumFastSync = 8000 // sync block num each goroutine
	ConsulAddr           = "127.0.0.1:8500"
	SyncWithDLock        = false
)

// get value of env var
func init() {
	nodeUrl, found := os.LookupEnv(constant.EnvNameSerNetworkNodeUrl)
	if found {
		BlockChainMonitorUrl = nodeUrl
	}
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameSerNetworkNodeUrl, BlockChainMonitorUrl)

	chainId, found := os.LookupEnv(constant.EnvNameSerNetworkChainId)
	if found {
		ChainId = chainId
	}
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameSerNetworkChainId, ChainId)

	token, found := os.LookupEnv(constant.EnvNameSerNetworkToken)
	if found {
		Token = token
	}
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameSerNetworkToken, Token)

	maxGoroutine, found := os.LookupEnv(constant.EnvNameSerMaxGoRoutine)
	if found {
		var err error
		SyncMaxGoroutine, err = strconv.Atoi(maxGoroutine)
		if err != nil {
			logger.Error.Fatalf("Convert str to int failed, env var is %v\n",
				constant.EnvNameSerMaxGoRoutine)
		}
	}
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameSerMaxGoRoutine, SyncMaxGoroutine)

	syncBlockNum, found := os.LookupEnv(constant.EnvNameSerSyncBlockNum)
	if found {
		var err error
		SyncBlockNumFastSync, err = strconv.Atoi(syncBlockNum)
		if err != nil {
			logger.Error.Fatalf("Convert str to int failed, env var is %v\n",
				constant.EnvNameSerSyncBlockNum)
		}
	}
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameSerSyncBlockNum, SyncBlockNumFastSync)

	consulAddr, found := os.LookupEnv(constant.EnvNameConsulAddr)
	if found {
		ConsulAddr = consulAddr
	}
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameConsulAddr, ConsulAddr)

	withDLock, found := os.LookupEnv(constant.EnvNameSyncWithDLock)
	if found {
		flag, err := strconv.ParseBool(withDLock)
		if err != nil {
			logger.Error.Fatalf("Convert str to bool failed, env var is %v\n",
				constant.EnvNameSyncWithDLock)
		}
		SyncWithDLock = flag
	}
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameSyncWithDLock, SyncWithDLock)
}
