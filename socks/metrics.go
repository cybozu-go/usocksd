package socks

import (
	"github.com/cybozu-go/usocksd/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	proxyElapsedHist = promauto.With(metrics.Registry).NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "elapsed",
		Help:      "provides the time elapsed, in seconds, between proxy start and end",
	}, []string{"result"})
	proxyRequestsInflightGauge = promauto.With(metrics.Registry).NewGauge(prometheus.GaugeOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "requests_inflight",
		Help:      "provides the number of requests currently in-flight",
	})

	socksResponseCounter = promauto.With(metrics.Registry).NewCounterVec(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: "socks",
		Name:      "response",
		Help:      "socks responses",
	}, []string{"version", "status"})
)
