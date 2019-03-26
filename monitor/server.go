package monitor

import (
	"fmt"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/monitor/status"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type Monitor struct {
	providers []MetricsProvider
}

func NewMonitor() *Monitor {
	var providers []MetricsProvider
	monitor := &Monitor{
		providers: providers,
	}
	monitor.AddMetricsProvider(status.PrometheusMetrics())
	return monitor
}

func (m *Monitor) AddMetricsProvider(provider MetricsProvider) *Monitor {
	m.providers = append(m.providers, provider)
	return m
}

func (m *Monitor) Start() {
	var startMetris = func() {
		for {
			t := time.NewTimer(time.Duration(5) * time.Second)
			select {
			case <-t.C:
				for _, provider := range m.providers {
					go provider.Report()
				}
			}
		}
	}

	go startMetris()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: promhttp.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err == nil {
			logger.Error("start monitor error", logger.String("error", err.Error()))
		}
	}()
}

type MetricsProvider interface {
	Report()
}
