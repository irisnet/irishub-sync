package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
)

func TestParseCoin(t *testing.T) {
	coin := ParseCoin("11111lc1")
	coinStr, _ := json.Marshal(coin)
	t.Log(string(coinStr))
}

func Test_getPrecision(t *testing.T) {
	exam1 := "4999999999999999999999"
	exam2 := "49999999999999999999.99"
	exam3 := "49999999.999999.99999"
	exam4 := "49999999.999999"
	exam5 := "4999999999999.9999999"

	data := getPrecision(exam1)
	if data == "4999999999999990000000" {
		t.Log("OK")
	} else {
		t.Fatal("Failed")
	}
	data = getPrecision(exam2)
	if data == "49999999999999900000" {
		t.Log("OK")

	} else {
		t.Fatal("Failed")
	}
	data = getPrecision(exam3)
	if data == "49999999.999999.99999" {
		t.Log("OK")
	} else {
		t.Fatal("Failed")
	}
	data = getPrecision(exam4)
	if data == "49999999.999999" {
		t.Log("OK")
	} else {
		t.Fatal("Failed")
	}
	data = getPrecision(exam5)
	fmt.Println(data)
	amt, _ := strconv.ParseFloat(data, 64)
	fmt.Println(amt)
	if data == "4999999999999.999" {
		t.Log("OK")
	} else {
		t.Fatal("Failed")
	}
}
