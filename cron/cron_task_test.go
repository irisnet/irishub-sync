package cron

import (
	"testing"
	"time"
	"github.com/irisnet/irishub-sync/store"
)

func TestUpdateUnknownOrEmptyTypeTxsByPage(t *testing.T) {
	store.Start()
	defer func() {
		store.Stop()
	}()
	num, err := UpdateUnknownOrEmptyTypeTxsByPage(0, 20)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(num)

}

func TestCronService_StartCronService(t *testing.T) {
	store.Start()
	defer func() {
		store.Stop()
	}()
	new(CronService).StartCronService()
	time.Sleep(1 * time.Minute)
}

