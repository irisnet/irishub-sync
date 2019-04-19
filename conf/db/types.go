package db

import (
	"os"

	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/util/constant"
)

var (
	Addrs    = "192.168.150.7:30000"
	User     = "iris"
	Passwd   = "irispassword"
	Database = "sync-iris"
)

// get value of env var
func init() {
	addrs, found := os.LookupEnv(constant.EnvNameDbAddr)
	if found {
		Addrs = addrs
	}
	logger.Info("Env Value", logger.String(constant.EnvNameDbAddr, Addrs))

	user, found := os.LookupEnv(constant.EnvNameDbUser)
	if found {
		User = user
	}

	passwd, found := os.LookupEnv(constant.EnvNameDbPassWd)
	if found {
		Passwd = passwd
	}

	database, found := os.LookupEnv(constant.EnvNameDbDataBase)
	if found {
		Database = database
	}
	logger.Info("Env Value", logger.String(constant.EnvNameDbDataBase, Database))
}
