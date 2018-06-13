package main

import (
	"sync"
	syncTask "github.com/irisnet/irishub-sync/sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	
	syncTask.Start()
	
	wg.Wait()
}
