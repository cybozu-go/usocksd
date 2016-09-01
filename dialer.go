package usocksd

import (
	"errors"
	"hash/fnv"
	"net"
	"strconv"
	"time"

	"github.com/cybozu-go/usocksd/socks"
)

const (
	dialTimeout = 10 * time.Second
)

type dialer struct {
	*AddressGroup
}

func calcHint(caddr, daddr net.IP) uint32 {
	hash := fnv.New32a()
	hash.Write(caddr)
	hash.Write(daddr)
	return hash.Sum32()
}

func (d dialer) Dial(r *socks.Request) (net.Conn, error) {
	var clientIP net.IP
	if tca, ok := r.Conn.RemoteAddr().(*net.TCPAddr); ok {
		clientIP = tca.IP
	}

	destIPs := []net.IP{r.IP}
	if len(r.Hostname) > 0 {
		ips, err := net.LookupIP(r.Hostname)
		if err != nil {
			return nil, err
		}
		destIPs = ips
	}

	deadline, ok := r.Context().Deadline()
	if !ok {
		deadline = time.Now().Add(dialTimeout)
	}

	var err error
	for _, ip := range destIPs {
		if time.Now().After(deadline) {
			err = errors.New("dial timeout")
			break
		}

		hint := calcHint(clientIP, ip)
		laddr := &net.TCPAddr{
			IP: d.PickAddress(hint),
		}
		raddr := &net.TCPAddr{
			IP:   ip,
			Port: r.Port,
		}
		conn, err2 := net.DialTCP("tcp", laddr, raddr)
		if err2 == nil {
			return conn, nil
		}
		err = err2
	}

	return nil, err
}

type dumbDialer struct {
	*net.Dialer
}

func (d dumbDialer) Dial(r *socks.Request) (net.Conn, error) {
	var addr string
	if len(r.Hostname) > 0 {
		addr = net.JoinHostPort(r.Hostname, strconv.Itoa(r.Port))
	} else {
		addr = net.JoinHostPort(r.IP.String(), strconv.Itoa(r.Port))
	}
	return d.DialContext(r.Context(), "tcp", addr)
}

func createDialer(c *Config) socks.Dialer {
	if len(c.Outgoing.Addresses) == 0 {
		return dumbDialer{
			&net.Dialer{
				KeepAlive: 3 * time.Minute,
				DualStack: true,
			},
		}
	}

	ag := NewAddressGroup(c.Outgoing.Addresses, c.Outgoing.DNSBLDomain)
	return dialer{ag}
}
