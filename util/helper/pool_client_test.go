// init client from clientPool.
// client is httpClient of tendermint

package helper

import (
	"testing"

	"github.com/irisnet/irishub-sync/module/logger"
)

func TestInitClientPool(t *testing.T) {
	a := []int{1, 2, 3}
	b := make([]int, 6, 6)
	for index, value := range a {
		b[index] = value
	}
	b[3] = 4
	logger.Info.Println(b)
}
