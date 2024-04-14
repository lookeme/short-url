package main

import (
	"flag"
	"github.com/lookeme/short-url/internal/configuration"
	"os"
)

var (
	networkCfg = configuration.NetworkCfg{}
	cfg        = configuration.Config{
		Network: &networkCfg,
	}
)

func parseFlags() {
	flag.StringVar(&networkCfg.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&networkCfg.BaseURL, "b", "http://localhost:8080", "base address")
	flag.Parse()
	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		networkCfg.ServerAddress = serverAddress
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		networkCfg.BaseURL = baseURL
	}
}
