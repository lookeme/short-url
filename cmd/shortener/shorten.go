package main

import (
	"context"
	"github.com/lookeme/short-url/internal/app/domain/user"
	"github.com/lookeme/short-url/internal/security"
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
	storage, err := createStorage(ctx, zlogger, cfg.Storage)
	if err != nil {
		return err
	}
	urlService := shorten.NewURLService(storage.ShortenRepository, zlogger, cfg)
	userService := user.NewUserService(storage.UserRepository, zlogger)
	urlHandler := handler.NewURLHandler(&urlService, &userService)
	authService := security.New(&userService, zlogger)
	var gzip compression.Compressor
	server := http.NewServer(urlHandler, cfg.Network, zlogger, &gzip, authService)
	defer storage.Close()
	return server.Serve()
}

func createStorage(ctx context.Context, log *logger.Logger, cfg *configuration.Storage) (*db.Storage, error) {
	var storage *db.Storage
	if len(cfg.ConnString) == 0 {
		shortenStore, err := inmemory.NewInMemShortenStorage(cfg, log)
		if err != nil {
			if err := shortenStore.RecoverFromFile(); err != nil {
				return storage, err
			}
		}
		userStore, err := inmemory.NewInMemUserStorage(log)
		if err != nil {
			return storage, err
		}
		storage = db.NewStorage(userStore, shortenStore)
	} else {
		postgres, err := db.New(ctx, log, cfg)
		if err != nil {
			return storage, err
		}
		shortenStorage := db.NewShortenRepository(postgres)
		userStorage := db.NewUserRepository(postgres)
		storage = db.NewStorage(userStorage, shortenStorage)
	}
	return storage, nil
}
