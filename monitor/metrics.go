package monitor

import (
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/helper"
	"net/http"
)

const (
	NodeStatusNotReachable = "not_reachable"
	NodeStatusCatchingUp   = "catching_up"
	NodeStatusSyncing      = "syncing"
)

func registerNetwork(r *mux.Router) {
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		var result struct {
			NodeHeight int64  `json:"node_height"`
			DbHeight   int64  `json:"db_height"`
			NodeStatus string `json:"node_status"`
		}
		client := helper.GetClient()
		defer func() {
			client.Release()
		}()
		status, err := client.Status()
		if err != nil {
			logger.Error("rpc node connection exception", logger.String("error", err.Error()))
			result.NodeStatus = NodeStatusNotReachable
			write(w, result)
			return
		}
		// node height
		result.NodeHeight = status.SyncInfo.LatestBlockHeight
		// db height
		result.DbHeight = document.Block{}.GetMaxBlockHeight()
		if status.SyncInfo.CatchingUp {
			result.NodeStatus = NodeStatusCatchingUp
		} else {
			result.NodeStatus = NodeStatusSyncing
		}
		write(w, result)
	}).Methods("GET")
}
