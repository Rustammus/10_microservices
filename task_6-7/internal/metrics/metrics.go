package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"microservices/task_6/internal/config"
	"net/http"
)

var (
	RequestsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_requests_total",
		Help: "The total number of requests",
	})
)

type AppMetrics struct {
	s http.Server
}

func NewAppMetrics() *AppMetrics {
	m := &AppMetrics{}
	m.s = http.Server{}
	go m.run()
	return m
}

func (m *AppMetrics) Close() {
	m.s.Shutdown(context.Background())
}

func (m *AppMetrics) RequestsIncrement() {
	RequestsProcessed.Inc()
}

func (m *AppMetrics) run() {

	c := config.GetConfig()

	m.s.Addr = ":" + c.Metrics.Port
	mux := http.NewServeMux()
	mux.Handle(c.Metrics.Path, promhttp.Handler())
	m.s.Handler = mux
	log.Print(m.s.ListenAndServe())
}
