package helper

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestLockAndUnLock(t *testing.T) {
	wc := sync.WaitGroup{}
	wc.Add(10)
	runner := func(flag int) {
		lock := NewLock("lock-key", 500*time.Millisecond)
		if lock.Lock() {
			fmt.Printf("Thread %d start! \n", flag)
			num := rand.Intn(2)
			time.Sleep(time.Duration(num) * time.Second)
			fmt.Printf("Thread %d end! \n", flag)
		}
		wc.Done()
		lock.UnLock()
	}

	go runner(1)
	go runner(2)
	go runner(3)
	go runner(4)
	go runner(5)
	go runner(6)
	go runner(7)
	go runner(8)
	go runner(9)
	go runner(10)
	wc.Wait()
}
