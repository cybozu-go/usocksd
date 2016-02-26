package main

import (
	"flag"
	"os"

	"github.com/cybozu-go/log"
	"github.com/cybozu-go/usocksd"
)

const (
	defaultConfigPath = "/usr/local/etc/usocksd.toml"
)

var (
	optFile = flag.String("f", "", "configuration file name")
)

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

	if len(c.Log.File) > 0 {
		mode := os.O_WRONLY | os.O_APPEND | os.O_CREATE
		f, err := os.OpenFile(c.Log.File, mode, 0644)
		if err != nil {
			log.ErrorExit(err)
		}
		defer f.Close()
		log.DefaultLogger().SetOutput(f)
	}
	err := log.DefaultLogger().SetThresholdByName(c.Log.Level)
	if err != nil {
		log.ErrorExit(err)
	}

	err = usocksd.ListenAndServe(c)
	log.Info("server ends", nil)
	if err != nil {
		log.ErrorExit(err)
	}
}
