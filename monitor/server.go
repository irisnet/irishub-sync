package monitor

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub-sync/logger"
	"net/http"
	"time"
)

type Monitor struct {
	*http.Server
}

func NewMonitor() *Monitor {
	router := getRouter()
	server := &http.Server{
		IdleTimeout: 10 * time.Second,
		Addr:        ":8080",
		Handler:     router,
	}
	return &Monitor{
		server,
	}
}

func (s *Monitor) Start() {
	logger.Info("#########################start monitor service##########################")
	go func() {
		if err := s.ListenAndServe(); err != nil {
			logger.Error("start monitor error", logger.String("error", err.Error()))
		}
	}()
}

func getRouter() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/monitor").Subrouter()
	registerNetwork(s)
	return r
}

func write(writer http.ResponseWriter, data interface{}) {
	if bz, err := json.Marshal(data); err == nil {
		writer.Write(bz)
	}
}
