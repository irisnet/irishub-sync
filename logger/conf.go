package logger

import (
	"github.com/irisnet/irishub-sync/util/constant"
	"os"
	"strconv"
)

type Config struct {
	Filename          string
	MaxSize           int
	MaxAge            int
	Compress          bool
	EnableAtomicLevel bool
}

var (
	conf = Config{
		Filename:          os.ExpandEnv("$HOME/irishub-sync/sync_server.log"),
		MaxSize:           20,
		MaxAge:            7,
		Compress:          true,
		EnableAtomicLevel: true,
	}
)

func init() {
	fileName, found := os.LookupEnv(constant.EnvLogFileName)
	if found {
		conf.Filename = fileName
	}

	maxSize, found := os.LookupEnv(constant.EnvLogFileMaxSize)
	if found {
		if size, err := strconv.Atoi(maxSize); err == nil {
			conf.MaxSize = size
		}
	}

	maxAge, found := os.LookupEnv(constant.EnvLogFileMaxAge)
	if found {
		if age, err := strconv.Atoi(maxAge); err == nil {
			conf.MaxAge = age
		}
	}

	compress, found := os.LookupEnv(constant.EnvLogCompress)
	if found {
		if compre, err := strconv.ParseBool(compress); err == nil {
			conf.Compress = compre
		}
	}

	enableAtomicLevel, found := os.LookupEnv(constant.EnableAtomicLevel)
	if found {
		if atomicLevel, err := strconv.ParseBool(enableAtomicLevel); err == nil {
			conf.EnableAtomicLevel = atomicLevel
		}
	}
}
