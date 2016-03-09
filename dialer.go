package usocksd

import (
	"net"
	"time"

	"golang.org/x/net/context"
)

const (
	dialTimeout       = 10 * time.Second
	keepAliveInterval = 300 * time.Second
)

func CreateDialer(c *Config) func(ctx context.Context, network, addr string) (net.Conn, error) {
	if len(c.Outgoing.Addresses) == 0 {
		return func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial(network, addr)
		}
	}

	ag := NewAddressGroup(c.Outgoing.Addresses, c.Outgoing.DNSBLDomain)
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		var hint uint32
		t := ctx.Value("hint")
		if t != nil {
			hint = t.(uint32)
		}

		d := net.Dialer{
			Timeout: dialTimeout,
			LocalAddr: &net.TCPAddr{
				IP: ag.PickAddress(hint),
			},
			KeepAlive: keepAliveInterval,
		}
		return d.Dial(network, addr)
	}
}
