package helper

import (
	"github.com/cosmos/cosmos-sdk/modules/coin"
)

var delay = false

func setDelay(d bool) {
	delay = d
}

func QueryAccountBalance(address string, delay bool) *coin.Account {
	// TODO: get balance of account from tendermint
	account := new(coin.Account)
	return account
}
