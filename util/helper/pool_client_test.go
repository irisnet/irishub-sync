// init client from clientPool.
// client is httpClient of tendermint

package helper

import (
	"testing"

	"github.com/tendermint/tendermint/rpc/client"

	conf "github.com/irisnet/irishub-sync/conf/server"
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

func TestGetClient(t *testing.T) {
	InitClientPool()

	for i := 0; i < conf.InitConnectionNum+10; i++ {
	}

}

func TestClient_Release(t *testing.T) {
	type fields struct {
		Client client.Client
		used   bool
		id     int64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Client{
				Client: tt.fields.Client,
				used:   tt.fields.used,
				id:     tt.fields.id,
			}
			n.Release()
		})
	}
}
