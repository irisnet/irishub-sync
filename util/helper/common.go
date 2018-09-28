package helper

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/module/logger"
	"strconv"
)

// convert object to json
func ToJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		logger.Error(err.Error())
	}
	return string(data)
}

func ParseStrToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
