package usocksd

import (
	"net"
	"time"
)

const (
	dialTimeout = 10 * time.Second
)

func CreateDialer(c *Config) func(network, addr string) (net.Conn, error) {
	if len(c.Outgoing.Addresses) == 0 {
		return net.Dial
	}

	ag := NewAddressGroup(c.Outgoing.Addresses, c.Outgoing.DNSBLDomain)
	return func(network, addr string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: dialTimeout,
			LocalAddr: &net.TCPAddr{
				IP: ag.PickAddress(),
			},
		}
		return d.Dial(network, addr)
	}
}
