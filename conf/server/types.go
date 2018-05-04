package server

const (
	BlockChainMonitorUrl = "tcp://localhost:46657"
	Token                = "iris"
	InitConnectionNum    = 100
	MaxConnectionNum     = 200
	ChainId              = "test"
	SyncCron             = "1 * * * * *"

	SyncMaxGoroutine     = 2000
	SyncBlockNumFastSync = 2000
)
