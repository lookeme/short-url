package main

import (
	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/compression"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/server/handler"
	"github.com/lookeme/short-url/internal/server/http"
	"github.com/lookeme/short-url/internal/storage/inmemory"
	"log"
)

func main() {
	cfg := configuration.CreateConfig()
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg *configuration.Config) error {
	zlogger, err := logger.CreateLogger(cfg.Logger)
	if err != nil {
		return err
	}
	storage := inmemory.NewStorage()
	urlService := shorten.NewURLService(storage)
	urlHandler := handler.NewURLHandler(urlService, cfg)
	gzip := &compression.Compressor{}
	server := http.NewServer(urlHandler, cfg.Network, zlogger, gzip)
	return server.Serve()
}
