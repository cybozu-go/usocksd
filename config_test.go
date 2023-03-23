package usocksd

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEmptyConfig(t *testing.T) {
	t.Parallel()

	c := NewConfig()
	expected := &Config{
		Incoming: IncomingConfig{
			Port:        defaultPort,
			MetricsPort: defaultMetricsPort,
		},
		Outgoing: OutgoingConfig{},
	}
	options := []cmp.Option{
		cmpopts.IgnoreUnexported(Config{}),
		cmpopts.IgnoreUnexported(IncomingConfig{}),
		cmpopts.IgnoreUnexported(OutgoingConfig{}),
	}
	if diff := cmp.Diff(c, expected, options...); diff != "" {
		t.Fatalf("unexpected config (-actual +expected):\n%s", diff)
	}
	if !c.allowPort(443) {
		t.Error("port 443 must be allowed")
	}
}

func TestConfig(t *testing.T) {
	t.Parallel()

	c := NewConfig()
	if err := c.Load("test/test1.toml"); err != nil {
		t.Fatal(err)
	}
	expected := &Config{
		Incoming: IncomingConfig{
			Port:        1080,
			MetricsPort: 8081,
			Addresses: []net.IP{
				net.ParseIP("127.0.0.1"),
			},
			AllowFrom: []string{
				"10.0.0.0/8",
				"192.168.1.1",
			},
		},
		Outgoing: OutgoingConfig{
			AllowSites: []string{
				"www.amazon.com",
				".google.com",
			},
			DenySites: []string{
				".2ch.net",
				"bad.google.com",
			},
			Addresses: []net.IP{
				net.ParseIP("12.34.56.78"),
			},
			DenyPorts:   []int{22, 25},
			DNSBLDomain: "zen.spamhaus.org",
		},
	}
	options := []cmp.Option{
		cmpopts.IgnoreUnexported(Config{}),
		cmpopts.IgnoreUnexported(IncomingConfig{}),
		cmpopts.IgnoreUnexported(OutgoingConfig{}),
	}
	if diff := cmp.Diff(c, expected, options...); diff != "" {
		t.Fatalf("unexpected config (-actual +expected):\n%s", diff)
	}
	if !c.allowIP(net.ParseIP("10.1.32.4")) {
		t.Error("10.1.32.4 is not allowed")
	}
	if !c.allowIP(net.ParseIP("192.168.1.1")) {
		t.Error("192.168.1.1 is not allowed")
	}
	if c.allowIP(net.ParseIP("12.34.56.78")) {
		t.Error("12.34.56.78 shout not be allowed")
	}
	if !c.allowFQDN("www.amazon.com") {
		t.Error("www.amazon.com should be allowed")
	}
	if !c.allowFQDN("www.google.com") {
		t.Error("www.google.com should be allowed")
	}
	if c.allowFQDN("bad.google.com") {
		t.Error("bad.amazon.com should be denied")
	}
	if c.allowFQDN("www.2ch.net") {
		t.Error("www.2ch.net should be denied")
	}
	if c.allowPort(25) {
		t.Error("port 25 must not be allowed")
	}
	if !c.allowPort(443) {
		t.Error("port 443 must be allowed")
	}
}

func TestConfigFail(t *testing.T) {
	t.Parallel()

	c := NewConfig()
	if err := c.Load("test/test2.toml"); err == nil {
		t.Error("loadConfig should fail for test2.toml")
	}

	// expect type mismatch
	c = NewConfig()
	if err := c.Load("test/test3.toml"); err == nil {
		t.Error("loadConfig should fail for test3.toml")
	}
}
