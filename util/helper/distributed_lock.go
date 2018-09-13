package helper

import (
	"github.com/hashicorp/consul/api"
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
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	return &DLock{
		Client:   client,
		lockKey:  lockKey,
		tryDelay: tryDelay,
	}
}

func (l *DLock) Lock() bool {
	session := l.Session()
	sessionId, _, err := session.CreateNoChecks(nil, nil)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	l.sessionId = sessionId
	p := &api.KVPair{Key: l.lockKey, Session: sessionId}
loop:
	{
		if work, _, err := l.KV().Acquire(p, nil); err != nil || !work {
			logger.Info.Printf("acquire lock[%s] fail,try again after 0.5s", l.lockKey)
			time.Sleep(l.tryDelay)
			goto loop
		}
		return true
	}
}

func (l *DLock) UnLock() {
	p := &api.KVPair{Key: l.lockKey, Session: l.sessionId}
	l.KV().Release(p, nil)
}

func (l *DLock) Destroy() {
	l.Session().Destroy(l.sessionId, nil)
}
