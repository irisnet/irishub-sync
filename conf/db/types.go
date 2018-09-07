package db

import (
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/util/constant"
	"os"
)

var (
	Host     = "192.168.150.7"
	Port     = "30000"
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
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameDbHost, Host)

	port, found := os.LookupEnv(constant.EnvNameDbPort)
	if found {
		Port = port
	}
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameDbPort, Port)

	user, found := os.LookupEnv(constant.EnvNameDbUser)
	if found {
		User = user
	}
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameDbUser, User)

	passwd, found := os.LookupEnv(constant.EnvNameDbPassWd)
	if found {
		Passwd = passwd
	}
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameDbPassWd, Passwd)

	database, found := os.LookupEnv(constant.EnvNameDbDataBase)
	if found {
		Database = database
	}
	logger.Info.Printf("The value of env var %v is %v\n",
		constant.EnvNameDbDataBase, Database)
}
