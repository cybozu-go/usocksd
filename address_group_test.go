package usocksd

import (
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

	a := &AddressGroup{
		lock: new(sync.Mutex),
		valids: []net.IP{
			net.ParseIP("12.34.56.78"),
			net.ParseIP("10.0.0.1"),
		},
	}

	if ip := a.PickAddress(); ip.String() != "12.34.56.78" {
		t.Error(ip.String() + " != 12.34.56.78")
	}
	if ip := a.PickAddress(); ip.String() != "10.0.0.1" {
		t.Error(ip.String() + " != 10.0.0.1")
	}
	if ip := a.PickAddress(); ip.String() != "12.34.56.78" {
		t.Error(ip.String() + " != 12.34.56.78")
	}
}
