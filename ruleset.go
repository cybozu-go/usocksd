package usocksd

import (
	"strconv"

	socks5 "github.com/cybozu-go/go-socks5"
	"github.com/cybozu-go/log"
	"golang.org/x/net/context"
)

type ruleSet struct {
	*Config
}

func (r ruleSet) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	if req.Command != socks5.ConnectCommand {
		return ctx, false
	}
	if !r.allowFQDN(req.DestAddr.FQDN) {
		log.Warn("denied access", map[string]interface{}{
			"_client_ip": req.RemoteAddr.IP.String(),
			"_fqdn":      req.DestAddr.FQDN,
		})
		return ctx, false
	}
	if !r.allowIP(req.RemoteAddr.IP) {
		log.Warn("denied access", map[string]interface{}{
			"_client_ip": req.RemoteAddr.IP.String(),
		})
		return ctx, false
	}
	if !r.allowPort(req.DestAddr.Port) {
		log.Warn("denied access", map[string]interface{}{
			"_client_ip": req.RemoteAddr.IP.String(),
			"_dest_port": strconv.Itoa(req.DestAddr.Port),
		})
		return ctx, false
	}
	return ctx, true
}

// CreateRuleSet returns a RuleSet for socks5.
func CreateRuleSet(c *Config) socks5.RuleSet {
	return ruleSet{c}
}
