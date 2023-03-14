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
		Name:      "elapsed_seconds",
		Help:      "provides the time elapsed, in seconds, between proxy start and end",
	}, []string{"result"})
	proxyRequestsInflightGauge = promauto.With(metrics.Registry).NewGauge(prometheus.GaugeOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "inflight_requests",
		Help:      "provides the number of requests currently in-flight",
	})
	proxyBytesTxHist = promauto.With(metrics.Registry).NewHistogram(prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "tx_bytes",
		Help:      "bytes copied from the source connection to the destination connection",
	})
	proxyBytesRxHist = promauto.With(metrics.Registry).NewHistogram(prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "rx_bytes",
		Help:      "bytes copied from the destination connection to the source connection",
	})

	proxyErrTxCount = promauto.With(metrics.Registry).NewCounter(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "tx_errors_total",
		Help:      "number of errors encountered copying from source to destination",
	})

	proxyErrRxCount = promauto.With(metrics.Registry).NewCounter(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "rx_errors_total",
		Help:      "number of errors encountered copying from destination to source",
	})

	proxyElapsedTxHist = promauto.With(metrics.Registry).NewHistogram(prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "tx_seconds",
		Help:      "time spent copying from source to destination",
	})

	proxyElapsedRxHist = promauto.With(metrics.Registry).NewHistogram(prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "rx_seconds",
		Help:      "time spent copying from destination to source",
	})

	authNegotiateCounter = promauto.With(metrics.Registry).NewCounterVec(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: "socks5",
		Name:      "auth_negotiated_total",
		Help:      "number of auth negotiation",
	}, []string{"type", "reason", "result"})

	addressReadCounter = promauto.With(metrics.Registry).NewCounterVec(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: "socks5",
		Name:      "address_read_total",
		Help:      "address read total count",
	}, []string{"reason", "result"})

	socksResponseCounter = promauto.With(metrics.Registry).NewCounterVec(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: "socks",
		Name:      "responses_total",
		Help:      "number of socks responses",
	}, []string{"version", "status", "reason"})
)
