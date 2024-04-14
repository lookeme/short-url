package main

import (
	"flag"
	"github.com/lookeme/short-url/internal/configuration"
)

var (
	networkCfg = configuration.NetworkCfg{}
	cfg        = configuration.Config{
		Network: &networkCfg,
	}
)

func parseFlags() {
	flag.StringVar(&networkCfg.StartAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&networkCfg.BaseAddress, "b", "http://localhost:8000", "base address")
	flag.Parse()
}
