// init client from clientPool.
// client is httpClient of tendermint

package helper

import (
	"errors"
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/robfig/cron"

	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/tendermint/tendermint/rpc/client"
)

var (
	pool = ClientPool{}
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

func (pool ClientPool) ping() {
	go func() {
		task := cron.New()
		task.AddFunc("0/5 * * * * *", func() {
			var clients []Client
			iterator := pool.iterator()
			for iterator.HasNext() {
				c := iterator.Next().(Client)
				if !c.heartbeat() {
					logger.Error.Printf("client node[%d] is stop", c.id)
					iterator.Remove()
				} else {
					clients = append(clients, c)
				}
			}
			pool.clients = clients
			logger.Info.Printf("current available client :%d\n", len(clients))
		})
		task.Start()
	}()
}

func (pool ClientPool) iterator() Iterator {
	var d []interface{}
	for _, data := range pool.clients {
		d = append(d, data)
	}
	return &ArrayIterator{
		data: &d,
	}
}

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
	pool.ping()
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
func (c Client) Release() {
	c.used = false
	pool.clients[c.id] = c
	pool.available++
	pool.used--
}

func (c Client) heartbeat() bool {
	logger.Info.Printf("client node[%d] heartbeat", c.id)
	if !c.used {
		http := c.Client.(*client.HTTP)
		if _, err := http.Health(); err != nil {
			return false
		}
	}
	return true
}

func createConnection(id int64) Client {
	tmClient := Client{
		Client: client.NewHTTP(conf.BlockChainMonitorUrl, "/websocket"),
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
			return tmClient, nil
		}
	}
	return Client{}, errors.New("pool is empty")
}
