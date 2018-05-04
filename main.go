package main

import (
	"fmt"

	"github.com/irisnet/iris-sync-server/module/logger"
)

func main() {
	i := 5
	j := 2
	logger.Info.Printf("result of j / i is %v", i/j)
	fmt.Println("This is test")
}
