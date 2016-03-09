package usocksd

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cybozu-go/log"
)

const (
	// invalidCheckInterval is the interval between checking DNSBL for
	// black-listed IP addresses.
	invalidCheckInterval = 15 * time.Second
)

// AddressGroup is a group of external IP addresses to be used
// for outgoing connections.  With the help of associated goroutines,
// IP addresses listed on DNSBL will be checked and excluded.
type AddressGroup struct {
	addresses   []net.IP // immutable
	dnsblDomain string

	lock     *sync.Mutex
	valids   []net.IP
	invalids []net.IP
}

func makeDNSBLDomain(domain string, ip net.IP) string {
	if len(domain) == 0 {
		return ""
	}
	ip = ip.To4()
	if ip == nil {
		return ""
	}
	return fmt.Sprintf("%d.%d.%d.%d.%s", ip[3], ip[2], ip[1], ip[0], domain)
}

// isBadIP returns true if IP is registered on DNSBL.
func (a *AddressGroup) isBadIP(ip net.IP) bool {
	d := makeDNSBLDomain(a.dnsblDomain, ip)
	if len(d) == 0 {
		return false
	}
	_, err := net.LookupIP(d)
	return err == nil
}

func toStringList(ips []net.IP) []string {
	sips := make([]string, 0, len(ips))
	for _, ip := range ips {
		sips = append(sips, ip.String())
	}
	return sips
}

// detectInvalid is a non-returning method, thus should be
// called as a goroutine, to detect black-listed IP addresses.
func (a *AddressGroup) detectInvalid() {
	for {
		var valids, invalids []net.IP
		for _, ip := range a.addresses {
			if a.isBadIP(ip) {
				invalids = append(invalids, ip)
			} else {
				valids = append(valids, ip)
			}
		}
		a.lock.Lock()
		if len(invalids) > 0 && len(a.invalids) != len(invalids) {
			log.Warn("detect black-listed IP", map[string]interface{}{
				"_bad_ips": toStringList(invalids),
			})
		}
		if len(valids) < len(invalids) {
			// Too few valid IPs
			valids = a.addresses
		}
		a.valids = valids
		a.invalids = invalids
		a.lock.Unlock()
		time.Sleep(invalidCheckInterval)
	}
}

// PickAddress returns a local IP address for outgoing connection.
// hint should be an integer calculated from client and/or target IP addresses.
func (a *AddressGroup) PickAddress(hint uint32) net.IP {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.valids[int(hint)%len(a.valids)]
}

// NewAddressGroup initializes a new AddressGroup and starts
// helper goroutines.
func NewAddressGroup(addresses []net.IP, dnsblDomain string) *AddressGroup {
	a := &AddressGroup{
		addresses:   addresses,
		dnsblDomain: dnsblDomain,
		lock:        new(sync.Mutex),
		valids:      addresses,
		invalids:    nil,
	}
	go a.detectInvalid()
	return a
}
