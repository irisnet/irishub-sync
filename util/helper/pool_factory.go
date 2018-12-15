package helper

import (
	"context"
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
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
	for _, url := range conf.BlockChainMonitorUrl {
		peersMap[generateId(url)] = EndPoint{
			Address:   url,
			Available: true,
		}
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

	logger.Info("PoolConfig", logger.Int("config.MaxTotal", config.MaxTotal), logger.Int("config.MaxIdle", config.MaxIdle))
	pool = &NodePool{
		gcp.NewObjectPool(ctx, &factory, config),
	}
	pool.PreparePool(ctx)
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
	logger.Info("release resource nodePool")
	pool.Close(ctx)
	factory.cron.Stop()
}

func (f *PoolFactory) MakeObject(ctx context.Context) (*gcp.PooledObject, error) {
	endpoint := f.GetEndPoint()
	logger.Debug("PoolFactory MakeObject peer", logger.Any("endpoint", endpoint))
	return gcp.NewPooledObject(newClient(endpoint.Address)), nil
}

func (f *PoolFactory) DestroyObject(ctx context.Context, object *gcp.PooledObject) error {
	logger.Debug("PoolFactory DestroyObject peer", logger.Any("peer", object.Object))
	c := object.Object.(*Client)
	if c.IsRunning() {
		c.Stop()
	}
	return nil
}

func (f *PoolFactory) ValidateObject(ctx context.Context, object *gcp.PooledObject) bool {
	// do validate
	logger.Debug("PoolFactory ValidateObject peer", logger.Any("peer", object.Object))
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
	logger.Debug("PoolFactory ActivateObject peer", logger.Any("peer", object.Object))
	return nil
}

func (f *PoolFactory) PassivateObject(ctx context.Context, object *gcp.PooledObject) error {
	logger.Debug("PoolFactory PassivateObject peer", logger.Any("peer", object.Object))
	return nil
}

func (f *PoolFactory) GetEndPoint() EndPoint {
	logger.Info("PoolFactory pullEndPoint peer", logger.Any("peers", f.peersMap))
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
			logger.Info("PoolFactory StartCrawlPeers peer", logger.Any("peers", f.peersMap))
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
