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
	"github.com/lookeme/short-url/internal/storage"
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
	store, err := createStorage(ctx, zlogger, cfg.Storage)
	if err != nil {
		return err
	}
	urlService := shorten.NewURLService(store, zlogger, cfg)
	urlHandler := handler.NewURLHandler(&urlService)
	var gzip compression.Compressor
	server := http.NewServer(urlHandler, cfg.Network, zlogger, &gzip)
	defer store.Close()
	return server.Serve()
}

func createStorage(ctx context.Context, log *logger.Logger, cfg *configuration.Storage) (storage.Repository, error) {
	if len(cfg.ConnString) == 0 {
		store, err := inmemory.NewStorage(cfg, log)
		if err != nil {
			if err := store.RecoverFromFile(); err != nil {
				return nil, err
			}
		}
		return store, err
	}
	return db.New(ctx, log, cfg)
}
