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
	storage, err := inmemory.NewStorage(cfg.Storage, zlogger)
	if err != nil {
		return err
	}
	if err := storage.RecoverFromFile(); err != nil {
		return err
	}
	urlService := shorten.NewURLService(storage, cfg)
	urlHandler := handler.NewURLHandler(urlService)
	gzip := &compression.Compressor{}
	server := http.NewServer(urlHandler, cfg.Network, zlogger, gzip)
	defer storage.Close()
	return server.Serve()
}
