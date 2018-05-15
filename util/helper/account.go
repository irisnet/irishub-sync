// This package is used for query balance of account

package helper

import (
	"fmt"
	"time"

	"github.com/irisnet/iris-sync-server/module/logger"

	wire "github.com/tendermint/go-wire"
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"github.com/cosmos/cosmos-sdk/stack"
	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/cosmos/cosmos-sdk/client/commands/query"
	"github.com/tendermint/go-wire/data"
	"github.com/tendermint/iavl"
	"github.com/cosmos/cosmos-sdk/client"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

var delay = false

func setDelay(d bool) {
	delay = d
}

func QueryAccountBalance(address string, delay bool) *coin.Account {
	account := new(coin.Account)
	actor, err := commands.ParseActor(address)
	if err != nil {
		return account
	}

	actor = coin.ChainAddr(actor)
	key := stack.PrefixedKey(coin.NameCoin, actor.Bytes())
	if delay {
		time.Sleep(1 * time.Second)
	}
	_, err2 := GetParsed(key, account, query.GetHeight(), false)
	if err2 != nil {
		logger.Info.Printf("account bytes are empty for address: %q\n", address)
	}
	return account
}


// argument (so pass in a pointer to the appropriate struct)
func GetParsed(key []byte, data interface{}, height int64, prove bool) (int64, error) {
	bs, h, err := Get(key, height, prove)
	if err != nil {
		return 0, err
	}
	err = wire.ReadBinaryBytes(bs, data)
	if err != nil {
		return 0, err
	}
	return h, nil
}

// Get queries the given key and returns the value stored there and the
// height we checked at.
//
// If prove is true (and why shouldn't it be?),
// the data is fully verified before returning.  If prove is false,
// we just repeat whatever any (potentially malicious) node gives us.
// Only use that if you are running the full node yourself,
// and it is localhost or you have a secure connection (not HTTP)
func Get(key []byte, height int64, prove bool) (data.Bytes, int64, error) {
	if height < 0 {
		return nil, 0, fmt.Errorf("Height cannot be negative\n")
	}

	if !prove {
		tmClient := GetClient()
		defer tmClient.Release()
		resp, err := tmClient.Client.ABCIQueryWithOptions("/key", key,
			rpcclient.ABCIQueryOptions{Trusted: true, Height: int64(height)})
		if resp == nil {
			return nil, height, err
		}
		return data.Bytes(resp.Response.Value), resp.Response.Height, err
	}
	val, h, _, err := GetWithProof(key, height)
	return val, h, err
}

// GetWithProof returns the values stored under a given key at the named
// height as in Get.  Additionally, it will return a validated merkle
// proof for the key-value pair if it exists, and all checks pass.
func GetWithProof(key []byte, height int64) (data.Bytes, int64, iavl.KeyProof, error) {
	tmClient := GetClient()
	defer tmClient.Release()
	cert, err := commands.GetCertifier()
	if err != nil {
		return nil, 0, nil, err
	}
	return client.GetWithProof(key, height, tmClient.Client, cert)
}

