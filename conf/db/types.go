package db

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/util/constant"
	"os"
)

var (
	Host     = "127.0.0.1"
	Port     = "27017"
	User     = "iris"
	Passwd   = "irispassword"
	Database = "sync-iris"
)

// get value of env var
func init() {
	host, found := os.LookupEnv(constant.EnvNameDbHost)
	if found {
		Host = host
	}
	logger.Info("Env Value", logger.String(constant.EnvNameDbHost, Host))

	port, found := os.LookupEnv(constant.EnvNameDbPort)
	if found {
		Port = port
	}
	logger.Info("Env Value", logger.String(constant.EnvNameDbPort, Port))

	user, found := os.LookupEnv(constant.EnvNameDbUser)
	if found {
		User = user
	}
	logger.Info("Env Value", logger.String(constant.EnvNameDbUser, User))

	passwd, found := os.LookupEnv(constant.EnvNameDbPassWd)
	if found {
		Passwd = passwd
	}
	logger.Info("Env Value", logger.String(constant.EnvNameDbPassWd, Passwd))

	database, found := os.LookupEnv(constant.EnvNameDbDataBase)
	if found {
		Database = database
	}
	logger.Info("Env Value", logger.String(constant.EnvNameDbDataBase, Database))
}
