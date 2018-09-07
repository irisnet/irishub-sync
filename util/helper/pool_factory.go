package helper

import (
	"context"
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/module/logger"
	gcp "github.com/jolestar/go-commons-pool"
	"github.com/robfig/cron"
)

var pool *gcp.ObjectPool
var ctx = context.Background()

func init() {
	peersMap := map[string]endPoint{}
	peersMap[generateId(conf.BlockChainMonitorUrl)] = endPoint{
		Address:   conf.BlockChainMonitorUrl,
		Available: true,
	}

	factory := PoolFactory{
		peersMap: peersMap,
	}
	config := gcp.NewDefaultPoolConfig()

	config.MaxTotal = conf.MaxConnectionNum
	config.MaxIdle = conf.InitConnectionNum
	config.MinIdle = conf.InitConnectionNum
	config.TestOnBorrow = true
	config.TestOnCreate = true
	config.TestWhileIdle = true

	logger.Info.Printf("MaxTotal %d ", config.MaxTotal)
	logger.Info.Printf("MaxIdle %d ", config.MaxIdle)
	logger.Info.Printf("MinIdle %d ", config.MinIdle)
	pool = gcp.NewObjectPool(ctx, &factory, config)
	//自动搜索可用节点
	factory.beginFetch()
}

type PoolFactory struct {
	peersMap map[string]endPoint
}

func (f *PoolFactory) MakeObject(ctx context.Context) (*gcp.PooledObject, error) {
	endpoint := f.pullEndPoint()
	logger.Info.Printf("PoolFactory MakeObject select peer[%v]  \n", endpoint)
	return gcp.NewPooledObject(newClient(endpoint.Address)), nil
}

func (f *PoolFactory) DestroyObject(ctx context.Context, object *gcp.PooledObject) error {
	c := object.Object.(*Client)
	return c.HeartBeat()
}

func (f *PoolFactory) ValidateObject(ctx context.Context, object *gcp.PooledObject) bool {
	// do validate
	c := object.Object.(*Client)
	if c.HeartBeat() != nil {
		if endPoint, ok := f.peersMap[c.Id]; ok {
			endPoint.Available = false
			f.peersMap[c.Id] = endPoint
		}
		return false
	}
	return true
}

func (f *PoolFactory) ActivateObject(ctx context.Context, object *gcp.PooledObject) error {
	// do activate
	c := object.Object.(*Client)
	err := c.HeartBeat()
	if err != nil {
		if endPoint, ok := f.peersMap[c.Id]; ok {
			logger.Info.Printf("PoolFactory ActivateObject peer[%s] is unavailable \n", endPoint.Address)
			endPoint.Available = false
			f.peersMap[c.Id] = endPoint
		}
	}
	return err
}

func (f *PoolFactory) PassivateObject(ctx context.Context, object *gcp.PooledObject) error {
	// do passivate
	c := object.Object.(*Client)
	err := c.HeartBeat()
	if err != nil {
		if endPoint, ok := f.peersMap[c.Id]; ok {
			logger.Info.Printf("PoolFactory peer[%s] is unavailable \n", endPoint.Address)
			endPoint.Available = false
			f.peersMap[c.Id] = endPoint
		}
	}
	return err
}

func (f *PoolFactory) pullEndPoint() endPoint {
	logger.Info.Printf("PoolFactory pullEndPoint peers %v ", f.peersMap)
	var key string
	for key = range f.peersMap {
		endPoint := f.peersMap[key]
		if endPoint.Available {
			return endPoint
		}
	}
	return f.peersMap[key]
}

func (f *PoolFactory) beginFetch() {
	go func() {
		c := cron.New()
		c.AddFunc("0 0/1 * * * *", func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Error.Printf("PoolFactory beginFetch error: %v ", err)
				}
			}()
			logger.Info.Printf("PoolFactory beginFetch peers %v ", f.peersMap)
			c, err := pool.BorrowObject(ctx)
			if err == nil {
				logger.Info.Printf("PoolFactory peers %v ########", c)
				http := c.(*Client)
				addrs := http.GetNodeAddress()
				for _, addr := range addrs {
					key := generateId(addr)
					if _, ok := f.peersMap[key]; !ok {
						f.peersMap[key] = endPoint{
							Address:   addr,
							Available: true,
						}
					}
				}
			}
			//检测节点是否上线
			for key := range f.peersMap {
				endPoint := f.peersMap[key]
				if !endPoint.Available {
					node := newClient(endPoint.Address)
					if node.HeartBeat() == nil {
						endPoint.Available = true
						f.peersMap[key] = endPoint
						node.Stop()
					}
				}
			}
		})
		c.Start()
	}()

}

type endPoint struct {
	Address   string
	Available bool
}
