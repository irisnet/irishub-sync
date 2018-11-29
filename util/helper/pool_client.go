// init client from clientPool.
// client is httpClient of tendermint

package helper

import (
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/types"
	"time"
	"fmt"
)

type Client struct {
	types.Client
	Id string
}

func newClient(addr string) *Client {
	return &Client{
		Client: types.NewHTTP(addr, "/websocket"),
		Id:     generateId(addr),
	}
}

// get client from pool
// while get a client from pool, available should -1, used should +1
func GetClient() *Client {

	c, err := pool.BorrowObject(ctx)
	for err != nil {
		logger.Error("GetClient failed,will try again after 3 seconds", logger.String("err", err.Error()))
		time.Sleep(3 * time.Second)
		c, err = pool.BorrowObject(ctx)
	}

	logger.Debug("current available connection", logger.Int("Num", pool.GetNumIdle()))
	logger.Debug("current used connection", logger.Int("Num", pool.GetNumActive()))
	return c.(*Client)
}

// release client
func (c *Client) Release() {
	err := pool.ReturnObject(ctx, c)
	if err != nil {
		logger.Debug("debug=======================Release err=======================debug")
		logger.Error(err.Error())
	}
	logger.Debug("debug=======================Release return=======================debug")
}

func (c *Client) HeartBeat() error {
	http := c.Client.(*types.HTTP)
	_, err := http.Health()
	return err
}


func (c *Client) GetNodeAddress() []string {
	//http := c.Client.(*types.HTTP)
	//netInfo, err := http.NetInfo()
	var addrs []string
	//if err == nil {
	//	peers := netInfo.Peers
	//	for _, peer := range peers {
	//		addr := peer.NodeInfo.ListenAddr
	//		ip := strings.Split(addr, ":")[0]
	//		port := strings.Split(peer.NodeInfo.Other[5], ":")[2] //TODO
	//		endpoint := fmt.Sprintf("%s%s:%s", "tcp://", ip, port)
	//		addrs = append(addrs, endpoint)
	//	}
	//}
	logger.Debug("found new node ", logger.Any("address", addrs))
	return addrs
}

func generateId(address string) string {
	return fmt.Sprintf("peer[%s]", address)
}
