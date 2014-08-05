package main

import (
	"log"
	"os"
	"strconv"

	"code.google.com/p/gcfg"
)

type Config struct {
	Server configServer
}

type configServer struct {
	Port     int
	Hostname string
}

func (c *configServer) Addr() string {
	return c.Hostname + ":" + strconv.Itoa(c.Port)
}

const defaultConfig = `
	[server]
	port=1987
	hostname=
`

func init() {
	logger = log.New(os.Stdout, "[websockify] ", log.Ldate|log.Ltime)

	err = gcfg.ReadStringInto(&cfg, defaultConfig)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
