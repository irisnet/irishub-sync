// init client from clientPool.
// client is httpClient of tendermint

package helper

import (
	"fmt"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/tendermint/tendermint/rpc/client"
	"time"
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
			logger.Error("GetClient err", logger.Any("err", err))
		}
	}()

	c, err := pool.BorrowObject(ctx)
	for err != nil {
		logger.Error("GetClient failed,will try again after 3 seconds", logger.String("err", err.Error()))
		time.Sleep(3 * time.Second)
		c, err = pool.BorrowObject(ctx)
	}
	logger.Info("current available connection", logger.Int("Num", pool.GetNumIdle()))
	logger.Info("current used connection", logger.Int("Num", pool.GetNumActive()))
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
	http := c.Client.(*client.HTTP)
	_, err := http.Health()
	return err
}

func generateId(address string) string {
	return fmt.Sprintf("peer[%s]", address)
}
