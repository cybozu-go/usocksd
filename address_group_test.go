package usocksd

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestMakeDNSBLDomain(t *testing.T) {
	t.Parallel()

	d := makeDNSBLDomain("zen.spamhaus.org", net.ParseIP("12.34.56.78"))
	if d != "78.56.34.12.zen.spamhaus.org" {
		t.Error(d + ` != "78.56.34.12.zen.spamhaus.org"`)
	}
}

func TestIsBadIP(t *testing.T) {
	t.Parallel()

	a := &AddressGroup{dnsblDomain: "zen.spamhaus.org"}

	// unstable test
	if false && !a.isBadIP(net.ParseIP("211.128.234.229")) {
		t.Error("211.128.234.229 should be black-listed")
	}

	if a.isBadIP(net.ParseIP("10.0.0.1")) {
		t.Error("10.0.0.1 should not be black-listed")
	}
}

func TestPick(t *testing.T) {
	t.Parallel()

	ip1 := net.ParseIP("12.34.56.78")
	ip2 := net.ParseIP("10.0.0.1")
	a := &AddressGroup{
		lock:   new(sync.Mutex),
		valids: []net.IP{ip1, ip2},
	}

	validate := func(ip net.IP) error {
		if !ip.Equal(ip1) && !ip.Equal(ip2) {
			return fmt.Errorf("invalid IP: %v", ip)
		}
		return nil
	}

	if err := validate(a.PickAddress(0)); err != nil {
		t.Error(err)
	}
	if err := validate(a.PickAddress(1)); err != nil {
		t.Error(err)
	}
	if err := validate(a.PickAddress(2)); err != nil {
		t.Error(err)
	}
}
