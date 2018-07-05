package helper

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/module/logger"
	"encoding/binary"
)


// convert object to json
func ToJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		logger.Error.Println(err)
	}
	return string(data)
}

// byte to int
func BytesToInt(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}
