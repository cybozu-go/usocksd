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
	proxyBytesTxHist = promauto.With(metrics.Registry).NewHistogram(prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "bytes_tx",
		Help:      "bytes copied from the source connection to the destination connection",
	})
	proxyBytesRxHist = promauto.With(metrics.Registry).NewHistogram(prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "bytes_rx",
		Help:      "bytes copied from the destination connection to the source connection",
	})
	proxyErrSrcCopyCount = promauto.With(metrics.Registry).NewCounter(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "error_copy_src",
		Help:      "number of errors encountered copying from source to destination",
	})
	proxyErrDestCopyCount = promauto.With(metrics.Registry).NewCounter(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: "proxy",
		Name:      "error_copy_dest",
		Help:      "number of errors encountered copying from destination to source",
	})

	socksResponseCounter = promauto.With(metrics.Registry).NewCounterVec(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: "socks",
		Name:      "response",
		Help:      "socks responses",
	}, []string{"version", "status"})
)
