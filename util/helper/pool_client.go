// init client from clientPool.
// client is httpClient of tendermint

package helper

import (
	"encoding/hex"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/types"
	"time"
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

func generateId(address string) string {
	id := []byte(address)
	return hex.EncodeToString(id)
}
