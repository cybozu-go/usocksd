package socks

import (
	"context"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/cybozu-go/log"
	"github.com/cybozu-go/netutil"
	"github.com/cybozu-go/well"
)

const (
	copyBufferSize     = 64 << 10
	negotiationTimeout = 10 * time.Second
)

var (
	dialer = &net.Dialer{
		DualStack: true,
	}
)

// Authenticator is the interface for user authentication.
// It should look Username and Password field in the request and
// returns true if authentication succeeds.
//
// Note that both Username and Password may be empty.
type Authenticator interface {
	Authenticate(r *Request) bool
}

// RuleSet is the interface for access control.
// It should look the request properties and returns true
// if the request matches rules.
type RuleSet interface {
	Match(r *Request) bool
}

// Dialer is the interface to establish connection to the destination.
type Dialer interface {
	Dial(r *Request) (net.Conn, error)
}

// Server implement SOCKS protocol.
type Server struct {
	// Auth can be used to authenticate a request.
	// If nil, all requests are allowed.
	Auth Authenticator

	// Rules can be used to test a request if it matches rules.
	// If nil, all requests passes.
	Rules RuleSet

	// Dialer is used to make connections to destination servers.
	// If nil, net.DialContext is used.
	Dialer Dialer

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
	server well.Server
	pool   *sync.Pool
}

func (s *Server) init() {
	if s.Logger == nil {
		s.Logger = log.DefaultLogger()
	}
	s.server.ShutdownTimeout = s.ShutdownTimeout
	s.server.Env = s.Env
	s.server.Handler = s.handleConnection
	s.pool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, copyBufferSize)
		},
	}
}

// Serve starts a goroutine to accept connections.
// This returns immediately.  l will be closed when s.Env is canceled.
// See https://godoc.org/github.com/cybozu-go/well#Server.Serve
func (s *Server) Serve(l net.Listener) {
	s.once.Do(s.init)
	s.server.Serve(l)
}

func (s *Server) dial(ctx context.Context, r *Request, network string) (net.Conn, error) {
	if s.Dialer != nil {
		return s.Dialer.Dial(r)
	}

	var addr string
	if len(r.Hostname) == 0 {
		addr = net.JoinHostPort(r.IP.String(), strconv.Itoa(r.Port))
	} else {
		addr = net.JoinHostPort(r.Hostname, strconv.Itoa(r.Port))
	}

	ctx, cancel := context.WithTimeout(ctx, negotiationTimeout)
	defer cancel()
	return dialer.DialContext(ctx, network, addr)
}

// handleConnection implements SOCKS protocol.
func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	conn.SetDeadline(time.Now().Add(negotiationTimeout))

	var preamble [2]byte
	_, err := io.ReadFull(conn, preamble[:])
	if err != nil {
		fields := well.FieldsFromContext(ctx)
		fields["client_addr"] = conn.RemoteAddr().String()
		fields[log.FnError] = err.Error()
		s.Logger.Error("failed to read preamble", fields)
		return
	}

	connVer := version(preamble[0])
	var destConn net.Conn
	switch connVer {
	case SOCKS4:
		destConn = s.handleSOCKS4(ctx, conn, preamble[1])
		if destConn == nil {
			return
		}
	case SOCKS5:
		destConn = s.handleSOCKS5(ctx, conn, preamble[1])
		if destConn == nil {
			return
		}
	default:
		fields := well.FieldsFromContext(ctx)
		fields["client_addr"] = conn.RemoteAddr().String()
		s.Logger.Error("unknown SOCKS version", fields)
		return
	}
	defer destConn.Close()
	netutil.SetKeepAlive(destConn)

	// negotiation completed.
	var zeroTime time.Time
	conn.SetDeadline(zeroTime)

	// do proxy
	st := time.Now()
	env := well.NewEnvironment(ctx)
	env.Go(func(ctx context.Context) error {
		buf := s.pool.Get().([]byte)
		_, err := io.CopyBuffer(destConn, conn, buf)
		s.pool.Put(buf)
		if hc, ok := destConn.(netutil.HalfCloser); ok {
			hc.CloseWrite()
		}
		if hc, ok := conn.(netutil.HalfCloser); ok {
			hc.CloseRead()
		}
		return err
	})
	env.Go(func(ctx context.Context) error {
		buf := s.pool.Get().([]byte)
		_, err := io.CopyBuffer(conn, destConn, buf)
		s.pool.Put(buf)
		if hc, ok := conn.(netutil.HalfCloser); ok {
			hc.CloseWrite()
		}
		if hc, ok := destConn.(netutil.HalfCloser); ok {
			hc.CloseRead()
		}
		return err
	})
	env.Stop()
	err = env.Wait()

	fields := well.FieldsFromContext(ctx)
	fields["elapsed"] = time.Since(st).Seconds()
	if err != nil {
		fields[log.FnError] = err.Error()
		s.Logger.Error("proxy ends with an error", fields)
		return
	}
	s.Logger.Info("proxy ends", fields)
}
