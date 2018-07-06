package main

import (
	"github.com/irisnet/irishub-sync/service"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	service.Start()

	wg.Wait()
}
