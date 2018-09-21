// init client from clientPool.
// client is httpClient of tendermint

package helper

import (
	"encoding/hex"
	"fmt"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/tendermint/tendermint/rpc/client"
	"strings"
)

type Client struct {
	client.Client
	Id string
}

func newClient(addr string) *Client {
	return &Client{
		Client: client.NewHTTP(addr, "/websocket"),
		Id:     generateId(addr),
	}
}

// get client from pool
// while get a client from pool, available should -1, used should +1
func GetClient() *Client {
	defer func() {
		if err := recover(); err != nil {
			logger.Error.Println(err)
		}
	}()
	c, err := pool.BorrowObject(ctx)
	if err != nil {
		logger.Error.Println("GetClient failed,err:", err)
		return nil
	}
	logger.Info.Printf("current available connection:%d", pool.GetNumIdle())
	logger.Info.Printf("current used connection:%d", pool.GetNumActive())
	return c.(*Client)
}

// release client
func (c *Client) Release() {
	err := pool.ReturnObject(ctx, c)
	if err != nil {
		logger.Info.Println("debug=======================Release err=======================debug")
		logger.Error.Println(err.Error())
	}
	logger.Info.Println("debug=======================Release return=======================debug")
}

func (c *Client) HeartBeat() error {
	http := c.Client.(*client.HTTP)
	_, err := http.Health()
	return err
}

func (c *Client) GetNodeAddress() []string {
	http := c.Client.(*client.HTTP)
	netInfo, err := http.NetInfo()
	var addrs []string
	if err == nil {
		peers := netInfo.Peers
		for _, peer := range peers {
			addr := peer.ListenAddr
			ip := strings.Split(addr, ":")[0]
			port := strings.Split(peer.Other[5], ":")[2]
			endpoint := fmt.Sprintf("%s%s:%s", "tcp://", ip, port)
			addrs = append(addrs, endpoint)
		}
	}
	fmt.Printf("#######################%v##################\n", addrs)
	return addrs
}

func generateId(address string) string {
	id := []byte(address)
	return hex.EncodeToString(id)
}
