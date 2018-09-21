package helper

import (
	"context"
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/module/logger"
	gcp "github.com/jolestar/go-commons-pool"
	"github.com/robfig/cron"
	"math/rand"
)

var (
	factory PoolFactory
	pool    *NodePool
	ctx     = context.Background()
)

func init() {
	peersMap := map[string]EndPoint{}
	peersMap[generateId(conf.BlockChainMonitorUrl)] = EndPoint{
		Address:   conf.BlockChainMonitorUrl,
		Available: true,
	}

	factory = PoolFactory{
		peersMap: peersMap,
		cron:     cron.New(),
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
	pool = &NodePool{
		gcp.NewObjectPool(ctx, &factory, config),
	}
	//自动搜索可用节点
	factory.StartCrawlPeers()
}

type EndPoint struct {
	Address   string
	Available bool
}

type NodePool struct {
	*gcp.ObjectPool
}

type PoolFactory struct {
	peersMap map[string]EndPoint
	cron     *cron.Cron
}

func ClosePool() {
	logger.Info.Printf("release resource :%s", "nodePool")
	pool.Close(ctx)
	factory.cron.Stop()
}

func (f *PoolFactory) MakeObject(ctx context.Context) (*gcp.PooledObject, error) {
	endpoint := f.GetEndPoint()
	logger.Info.Printf("PoolFactory MakeObject peer[%v]  \n", endpoint)
	return gcp.NewPooledObject(newClient(endpoint.Address)), nil
}

func (f *PoolFactory) DestroyObject(ctx context.Context, object *gcp.PooledObject) error {
	logger.Info.Printf("PoolFactory DestroyObject peer[%v] \n", object.Object)
	c := object.Object.(*Client)
	if c.IsRunning() {
		c.Stop()
	}
	return nil
}

func (f *PoolFactory) ValidateObject(ctx context.Context, object *gcp.PooledObject) bool {
	// do validate
	logger.Info.Printf("PoolFactory ValidateObject peer[%v] \n", object.Object)
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
	logger.Info.Printf("PoolFactory ActivateObject peer[%v] \n", object.Object)
	return nil
}

func (f *PoolFactory) PassivateObject(ctx context.Context, object *gcp.PooledObject) error {
	logger.Info.Printf("PoolFactory PassivateObject peer[%v] \n", object.Object)
	return nil
}

func (f *PoolFactory) GetEndPoint() EndPoint {
	logger.Info.Printf("PoolFactory pullEndPoint peers %v ", f.peersMap)
	var keys []string
	var selectedKey string
	for key := range f.peersMap {
		endPoint := f.peersMap[key]
		if endPoint.Available {
			keys = append(keys, key)
		}
		selectedKey = key
	}
	if len(keys) > 0 {
		index := rand.Intn(len(keys))
		selectedKey = keys[index]
	}
	return f.peersMap[selectedKey]
}

func (f *PoolFactory) StartCrawlPeers() {
	go func() {
		f.cron.AddFunc("0 0/1 * * * *", func() {
			logger.Info.Printf("PoolFactory StartCrawlPeers peers %v ", f.peersMap)
			client := GetClient()
			logger.Info.Printf("PoolFactory peers %v ", client)

			defer func() {
				client.Release()
				if err := recover(); err != nil {
					logger.Error.Printf("PoolFactory StartCrawlPeers error: %v ", err)
				}
			}()

			addrs := client.GetNodeAddress()
			for _, addr := range addrs {
				key := generateId(addr)
				if _, ok := f.peersMap[key]; !ok {
					f.peersMap[key] = EndPoint{
						Address:   addr,
						Available: true,
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
					}
				}
			}
		})
		f.cron.Start()
	}()
}
