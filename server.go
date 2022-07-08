package usocksd

import (
	"net"
	"strconv"

	"github.com/cybozu-go/usocksd/metrics"
	"github.com/cybozu-go/usocksd/socks"
)

// Listeners returns a list of net.Listener.
func Listeners(c *Config) ([]net.Listener, error) {
	if len(c.Incoming.Addresses) == 0 {
		ln, err := net.Listen("tcp", ":"+strconv.Itoa(c.Incoming.Port))
		if err != nil {
			return nil, err
		}
		return []net.Listener{ln}, nil
	}

	lns := make([]net.Listener, len(c.Incoming.Addresses))
	for i, a := range c.Incoming.Addresses {
		addr := net.JoinHostPort(a.String(), strconv.Itoa(c.Incoming.Port))
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			for j := 0; j < i; j++ {
				lns[j].Close()
			}
			return nil, err
		}
		lns[i] = ln
	}
	return lns, nil
}

// MetricsListener returns a listener for the metrics server.
func MetricsListener(c *Config) (net.Listener, error) {
	addr := net.JoinHostPort("0.0.0.0", strconv.Itoa(c.Incoming.MetricsPort))
	return net.Listen("tcp", addr)
}

// NewServer creates a new socks.Server.
func NewServer(c *Config) *socks.Server {
	return &socks.Server{
		Rules:  createRuleSet(c),
		Dialer: createDialer(c),
	}
}

// NewMetricsServer creates a new metrics.Server.
func NewMetricsServer(c *Config) *metrics.Server {
	return &metrics.Server{}
}
