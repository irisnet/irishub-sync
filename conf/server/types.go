package server

const (
	BlockChainMonitorUrl = "tcp://47.104.155.125:46757"
	Token                = "iris"
	InitConnectionNum    = 1000
	MaxConnectionNum     = 2000
	ChainId              = "test"
	SyncCron             = "1 * * * * *"

	SyncMaxGoroutine     = 2000
	SyncBlockNumFastSync = 2000
)
