package usocksd

import (
	"errors"
	"net"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	DefaultLogLevel = "info"
	DefaultPort     = 1080
)

type LogConfig struct {
	File  string
	Level string
}

type IncomingConfig struct {
	Port         int
	Addresses    []net.IP
	AllowFrom    []string `toml:"allow_from"`
	allowSubnets []*net.IPNet
}

type OutgoingConfig struct {
	AllowSites  []string `toml:"allow_sites"`
	DenySites   []string `toml:"deny_sites"`
	DenyPorts   []int    `toml:"deny_ports"`
	Addresses   []net.IP
	DNSBLDomain string `toml:"dnsbl_domain"`
}

type Config struct {
	Log      LogConfig
	Incoming IncomingConfig
	Outgoing OutgoingConfig
}

func NewConfig() *Config {
	c := new(Config)
	c.Log.Level = DefaultLogLevel
	c.Incoming.Port = DefaultPort
	return c
}

func (c *Config) Load(path string) error {
	if md, err := toml.DecodeFile(path, c); err != nil {
		return err
	} else {
		if len(md.Undecoded()) > 0 {
			return errors.New("Unknown config keys in " + path)
		}
	}

	if len(c.Incoming.AllowFrom) > 0 {
		subnets := make([]*net.IPNet, 0, len(c.Incoming.AllowFrom))
		for _, s := range c.Incoming.AllowFrom {
			if strings.IndexByte(s, '/') == -1 {
				s = s + "/32"
			}
			_, n, err := net.ParseCIDR(s)
			if err != nil {
				return errors.New("Invalid network or IP address: " + s)
			}
			subnets = append(subnets, n)
		}
		c.Incoming.allowSubnets = subnets
	}

	c.Outgoing.AllowSites = toLowerStrings(c.Outgoing.AllowSites)
	c.Outgoing.DenySites = toLowerStrings(c.Outgoing.DenySites)

	return nil
}

func toLowerStrings(l []string) (nl []string) {
	for _, s := range l {
		nl = append(nl, strings.ToLower(s))
	}
	return
}

// allowIP tests if ip is allowed to connect to usocksd.
func (c *Config) allowIP(ip net.IP) bool {
	if len(c.Incoming.allowSubnets) == 0 {
		return true
	}
	for _, n := range c.Incoming.allowSubnets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

func siteMatch(site, match string) bool {
	if len(match) > 0 && match[0] == '.' {
		return strings.HasSuffix(site, match)
	}
	return site == match
}

// allowFQDN tests if FQDN is granted to access or not.
func (c *Config) allowFQDN(fqdn string) bool {
	fqdn = strings.ToLower(fqdn)
	if len(c.Outgoing.AllowSites) > 0 {
		for _, match := range c.Outgoing.AllowSites {
			if siteMatch(fqdn, match) {
				goto CHECK_DENY
			}
		}
		return false
	}

CHECK_DENY:
	for _, match := range c.Outgoing.DenySites {
		if siteMatch(fqdn, match) {
			return false
		}
	}
	return true
}

// allowPort tests if port is legitimate for destination.
func (c *Config) allowPort(port int) bool {
	for _, p := range c.Outgoing.DenyPorts {
		if p == port {
			return false
		}
	}
	return true
}
