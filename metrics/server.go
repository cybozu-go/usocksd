package metrics

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// Namespace provides a common namespace for metrics.
	Namespace = "usocksd"
)

var (
	// Registry provides a common registry for metrics.
	Registry       = prometheus.NewRegistry()
	metricsHandler = promhttp.HandlerFor(Registry, promhttp.HandlerOpts{
		ErrorHandling:     promhttp.ContinueOnError,
		EnableOpenMetrics: true,
	})
)

// Server implements a metrics server.
type Server struct {
	// Logger can be used to provide a custom logger.
	// If nil, the default logger is used.
	Logger *log.Logger

	// ShutdownTimeout is the maximum duration the server waits for
	// all connections to be closed before shutdown.
	//
	// Zero duration disables timeout.
	ShutdownTimeout time.Duration

	// Env is the environment where this server runs.
	//
	// The global environment is used if Env is nil.
	Env *well.Environment

	once   sync.Once
	server well.HTTPServer
}

func (s *Server) init() {
	if s.Logger == nil {
		s.Logger = log.DefaultLogger()
	}
	s.server.ShutdownTimeout = s.ShutdownTimeout
	s.server.Env = s.Env
	s.server.Server = &http.Server{
		Handler: metricsHandler,
	}
}

// Serve starts a goroutine to accept connections.
// This returns immediately.  l will be closed when s.Env is canceled.
// See https://godoc.org/github.com/cybozu-go/well#Server.Serve
func (s *Server) Serve(ln net.Listener) {
	s.once.Do(s.init)
	_ = s.server.Serve(ln)
}
