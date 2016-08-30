package usocksd

import (
	"net"

	"github.com/cybozu-go/log"
	"github.com/cybozu-go/usocksd/socks"
)

type ruleSet struct {
	*Config
}

func (ru ruleSet) Match(r *socks.Request) bool {
	clientAddr := r.Conn.RemoteAddr()
	tca, ok := clientAddr.(*net.TCPAddr)

	if !ru.allowFQDN(r.Hostname) {
		log.Warn("denied access", map[string]interface{}{
			"client_addr": clientAddr.String(),
			"fqdn":        r.Hostname,
		})
		return false
	}

	if ok && !ru.allowIP(tca.IP) {
		log.Warn("denied access", map[string]interface{}{
			"client_addr": clientAddr.String(),
		})
		return false
	}

	if !ru.allowPort(r.Port) {
		log.Warn("denied access", map[string]interface{}{
			"client_addr": clientAddr.String(),
			"dest_port":   r.Port,
		})
		return false
	}

	return true
}

func createRuleSet(c *Config) socks.RuleSet {
	return ruleSet{c}
}
