package helper

import (
	"encoding/json"
	"fmt"
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

func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.Error("common.ParseFloat error", logger.String("value", s))
	}
	return f
}

func RoundFloat(num float64, bit int) (i float64) {
	format := "%" + fmt.Sprintf("0.%d", bit) + "f"
	s := fmt.Sprintf(format, num)
	i, err := strconv.ParseFloat(s, 0)
	if err != nil {
		logger.Error("common.RoundFloat error", logger.String("format", format))
	}
	return i
}
