package usocksd

import (
	"fmt"
	_log "log"

	socks5 "github.com/armon/go-socks5"
	"github.com/cybozu-go/log"
)

func ListenAndServe(c *Config) error {
	logger := _log.New(log.DefaultLogger().Writer(log.LvError), "", 0)

	socksConfig := &socks5.Config{
		Rules:  CreateRuleSet(c),
		Logger: logger,
		Dial:   CreateDialer(c),
	}
	s, err := socks5.New(socksConfig)
	if err != nil {
		return err
	}

	log.Info("server starts", nil)

	if len(c.Incoming.Addresses) == 0 {
		return s.ListenAndServe("tcp", fmt.Sprintf(":%d", c.Incoming.Port))
	}

	ch := make(chan error)

	for _, ip := range c.Incoming.Addresses {
		addr := fmt.Sprintf("%v:%d", ip, c.Incoming.Port)
		go func() {
			ch <- s.ListenAndServe("tcp", addr)
		}()
	}

	return <-ch
}
