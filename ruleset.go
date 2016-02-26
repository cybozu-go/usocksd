package usocksd

import (
	"strconv"

	socks5 "github.com/armon/go-socks5"
	"github.com/cybozu-go/log"
)

type ruleSet struct {
	*Config
}

func (r ruleSet) Allow(req *socks5.Request) bool {
	if req.Command != socks5.ConnectCommand {
		return false
	}
	if !r.allowFQDN(req.DestAddr.FQDN) {
		log.Warn("denied access", map[string]interface{}{
			"_client_ip": req.RemoteAddr.IP.String(),
			"_fqdn":      req.DestAddr.FQDN,
		})
		return false
	}
	if !r.allowIP(req.RemoteAddr.IP) {
		log.Warn("denied access", map[string]interface{}{
			"_client_ip": req.RemoteAddr.IP.String(),
		})
		return false
	}
	if !r.allowPort(req.DestAddr.Port) {
		log.Warn("denied access", map[string]interface{}{
			"_client_ip": req.RemoteAddr.IP.String(),
			"_dest_port": strconv.Itoa(req.DestAddr.Port),
		})
		return false
	}
	return true
}

// CreateRuleSet returns a RuleSet for socks5.
func CreateRuleSet(c *Config) socks5.RuleSet {
	return ruleSet{c}
}
