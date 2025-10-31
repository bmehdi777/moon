package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	TunActiveConnection prometheus.Gauge
	TunHttpRequestTotal prometheus.Counter
	ApiHttpRequestTotal prometheus.Counter
}

func newMetrics() *metrics {
	m := metrics{
		TunActiveConnection: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "tun_active_connection",
		}),
		TunHttpRequestTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "tun_http_request_total",
		}),
		ApiHttpRequestTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "api_http_request_total",
		}),
	}

	defaultRegistry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = defaultRegistry
	prometheus.DefaultGatherer = defaultRegistry

	prometheus.MustRegister(m.TunActiveConnection, m.TunHttpRequestTotal, m.ApiHttpRequestTotal)

	return &m
}

var Metrics = newMetrics()
