package helper

import (
	"github.com/hashicorp/consul/api"
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/module/logger"
	"log"
	"time"
)

type DLock struct {
	*api.Client
	sessionId string
	lockKey   string
	tryDelay  time.Duration
}

func NewLock(lockKey string, tryDelay time.Duration) *DLock {
	config := api.DefaultConfig()
	config.Address = conf.ConsulAddr
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	session := client.Session()
	sessionId, _, err := session.CreateNoChecks(nil, nil)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	return &DLock{
		Client:    client,
		sessionId: sessionId,
		lockKey:   lockKey,
		tryDelay:  tryDelay,
	}
}

func (l *DLock) Lock() bool {
	p := &api.KVPair{Key: l.lockKey, Session: l.sessionId}
loop:
	{
		if work, _, err := l.KV().Acquire(p, nil); err != nil || !work {
			logger.Info("acquire lock fail", logger.String("lockKey", l.lockKey), logger.Duration("tryDelay", l.tryDelay))
			time.Sleep(l.tryDelay)
			goto loop
		}
		return true
	}
}

func (l *DLock) UnLock() {
	defer func() {
		//防止异常，锁无法释放
		if r := recover(); r != nil {
			l.Destroy()
		}
	}()
	p := &api.KVPair{Key: l.lockKey, Session: l.sessionId}
	l.KV().Release(p, nil)
}

func (l *DLock) Destroy() {
	l.Session().Destroy(l.sessionId, nil)
}
