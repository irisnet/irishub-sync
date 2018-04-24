// init client from clientPool.
// client is httpClient of tendermint

package helper

import (
	"errors"
	"log"

	conf "github.com/irisnet/iris-sync-server/conf/server"

	rpcClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/tendermint/tendermint/rpc/client"
)

type Client struct {
	Client client.Client
	used   bool
	id     int64
}

type ClientPool struct {
	clients        []Client
	available      int64
	used           int64
	maxConnection  int64
	initConnection int64
}

var pool = ClientPool{}

// init clientPool
// while init a client, available which is a var of clientPool should +1
func InitClientPool() {
	pool.maxConnection = int64(conf.MaxConnectionNum)
	initConnectionNum := int64(conf.InitConnectionNum)
	pool.initConnection = initConnectionNum

	pool.clients = make([]Client, initConnectionNum)
	for i := int64(0); i < initConnectionNum; i++ {
		createConnection(i)
	}
}

// get client from pool
// while get a client from pool, available should -1, used should +1
func GetClient() Client {
	c, err := getClient()
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// release client
func (n Client) Release() {
	n.used = false
	pool.clients[n.id] = n
	pool.available++
	pool.used--
}

func createConnection(id int64) Client {
	client := Client{
		Client: rpcClient.GetNode(conf.BlockChainMonitorUrl),
		used:   false,
		id:     id,
	}
	pool.clients[id] = client
	pool.available++
	return client
}

func getClient() (Client, error) {
	if pool.available == 0 {
		maxConnNum := int64(conf.MaxConnectionNum)
		if pool.used < maxConnNum {
			var client Client
			for i := int64(len(pool.clients)); i < maxConnNum; i++ {
				client = createConnection(i)
			}
			return client, nil
		} else {
			log.Fatal("client pool has no available connection")
		}
	}

	for _, client := range pool.clients {
		if !client.used {
			client.used = true
			pool.clients[client.id] = client
			pool.available--
			pool.used++
			log.Printf("current available coonection ：%d", pool.available)
			return client, nil
		}
	}
	return Client{}, errors.New("pool is empty")
}
