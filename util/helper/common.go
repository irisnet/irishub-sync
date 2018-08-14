package helper

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/module/logger"
)

// convert object to json
func ToJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		logger.Error.Println(err)
	}
	return string(data)
}
