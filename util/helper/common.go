package helper

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-sync/logger"
	"reflect"
	"strconv"
	"strings"
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

func ParseFloat(s string, bit ...int) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.Error("common.ParseFloat error", logger.String("value", s))
		return 0
	}

	if len(bit) > 0 {
		return RoundFloat(f, bit[0])
	}
	return f
}

func RoundFloat(num float64, bit int) (i float64) {
	format := "%" + fmt.Sprintf("0.%d", bit) + "f"
	s := fmt.Sprintf(format, num)
	i, err := strconv.ParseFloat(s, 0)
	if err != nil {
		logger.Error("common.RoundFloat error", logger.String("format", format))
		return 0
	}
	return i
}

func RoundString(decimal string, bit int) (i string) {
	f := ParseFloat(decimal, bit)
	return strconv.FormatFloat(f, 'f', bit, 64)
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		key := tag.Get("json")
		if len(key) == 0 {
			key = strings.ToLower(t.Field(i).Name)
		}
		data[key] = v.Field(i).Interface()
	}
	return data
}

func Map2Struct(srcMap map[string]interface{}, obj interface{}) {
	bz, err := json.Marshal(srcMap)
	if err != nil {
		logger.Error("map convert to struct failed")
	}
	err = json.Unmarshal(bz, obj)
	if err != nil {
		logger.Error("map convert to struct failed")
	}
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func DistinctStringSlice(slice []string) []string {
	var res []string
	elementExistMap := make(map[string]bool)
	if len(slice) > 0 {
		for _, v := range slice {
			if !elementExistMap[v] {
				res = append(res, v)
				elementExistMap[v] = true
			}
		}
	}

	return res
}
