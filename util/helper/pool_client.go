// init client from clientPool.
// client is httpClient of tendermint

package helper

import (
	"errors"
	conf "github.com/irisnet/iris-sync-server/conf/server"

	rpcClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/irisnet/iris-sync-server/module/logger"
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
		logger.Error.Fatalln(err.Error())
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
	tmClient := Client{
		Client: rpcClient.GetNode(conf.BlockChainMonitorUrl),
		used:   false,
		id:     id,
	}
	
	if id == int64(len(pool.clients)) {
		newSlice := make([]Client, pool.maxConnection)
		for i, v := range pool.clients {
			newSlice[i] = v
		}
		pool.clients = newSlice
	}
	pool.clients[id] = tmClient
	pool.available++
	return tmClient
}

func getClient() (Client, error) {
	if pool.available == 0 {
		maxConnNum := int64(conf.MaxConnectionNum)
		if pool.used < maxConnNum {
			var tmClient Client
			length := len(pool.clients)
			for i := int64(length); i < maxConnNum; i++ {
				tmClient = createConnection(i)
			}
			return tmClient, nil
		} else {
			logger.Error.Fatalln("client pool has no available connection")
		}
	}

	for _, tmClient := range pool.clients {
		if !tmClient.used {
			tmClient.used = true
			pool.clients[tmClient.id] = tmClient
			pool.available--
			pool.used++
			logger.Info.Printf("current available coonection ：%d\n", pool.available)
			return tmClient, nil
		}
	}
	return Client{}, errors.New("pool is empty")
}
