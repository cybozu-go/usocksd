package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/cybozu-go/log"
	"github.com/cybozu-go/usocksd"
	"github.com/cybozu-go/well"
)

const (
	defaultConfigPath = "/etc/usocksd.toml"
)

var (
	optFile = flag.String("f", "", "configuration file name")
)

func serveMetrics(c *usocksd.Config) error {
	metricsServer := usocksd.NewMetricsServer(c)
	mln, err := usocksd.MetricsListener(c)
	if err != nil {
		return fmt.Errorf("could not initialize metrics server: %w", err)
	}
	metricsServer.Serve(mln)
	return nil
}

func serve(lns []net.Listener, c *usocksd.Config) {
	socksServer := usocksd.NewServer(c)
	for _, ln := range lns {
		socksServer.Serve(ln)
	}
	if err := serveMetrics(c); err != nil {
		log.ErrorExit(err)
	}
	if err := well.Wait(); err != nil && !well.IsSignaled(err) {
		log.ErrorExit(err)
	}
}

func main() {
	flag.Parse()

	c := usocksd.NewConfig()
	if len(*optFile) > 0 {
		if err := c.Load(*optFile); err != nil {
			log.ErrorExit(err)
		}
	} else {
		_, err := os.Stat(defaultConfigPath)
		if err == nil {
			if e := c.Load(defaultConfigPath); e != nil {
				log.ErrorExit(e)
			}
		}
	}
	err := c.Log.Apply()
	if err != nil {
		log.ErrorExit(err)
	}

	g := &well.Graceful{
		Listen: func() ([]net.Listener, error) {
			return usocksd.Listeners(c)
		},
		Serve: func(lns []net.Listener) {
			serve(lns, c)
		},
	}
	g.Run()

	err = well.Wait()
	if err != nil && !well.IsSignaled(err) {
		log.ErrorExit(err)
	}
}
