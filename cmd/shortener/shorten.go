package main

import (
	"context"
	"log"

	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/compression"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/server/handler"
	"github.com/lookeme/short-url/internal/server/http"
	"github.com/lookeme/short-url/internal/storage/db"
	"github.com/lookeme/short-url/internal/storage/inmemory"
)

func main() {
	cfg := configuration.New()
	ctx := context.Background()
	if err := run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cfg *configuration.Config) error {
	zlogger, err := logger.CreateLogger(cfg.Logger)
	if err != nil {
		return err
	}
	storage, err := inmemory.NewStorage(cfg.Storage, zlogger)
	if err != nil {
		return err
	}
	dbStorage, err := db.NewDBStorage(ctx, zlogger, cfg.Storage)
	if err != nil {
		return err
	}
	if err := storage.RecoverFromFile(); err != nil {
		return err
	}
	urlService := shorten.NewURLService(storage, dbStorage, cfg)
	urlHandler := handler.NewURLHandler(&urlService)
	var gzip compression.Compressor
	server := http.NewServer(urlHandler, cfg.Network, zlogger, &gzip)
	defer storage.Close()
	return server.Serve()
}
