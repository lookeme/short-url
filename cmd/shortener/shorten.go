// Package main implements a URL shortener application.
//
// The application shortens URLs by generating a hash
// of the original URL, then storing it in a database.
// When the short representation of the URL is accessed,
// it redirects to the original URL.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lookeme/short-url/internal/app/domain/user"
	"github.com/lookeme/short-url/internal/security"
	"go.uber.org/zap"

	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/compression"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/server/handler"
	"github.com/lookeme/short-url/internal/server/http"
	"github.com/lookeme/short-url/internal/storage/db"
	"github.com/lookeme/short-url/internal/storage/inmemory"
)

// buildVersion represents the version number of the software build.
// go run -ldflags "-X main.buildVersion=v1.0.0" main.go
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// main is the entry point of the URL shortening application.
// It sets up the necessary configuration and starts the application.
// If something fails during setup or execution, the program will log a Fatal error message.
func main() {
	cfg := configuration.New()
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	run(cfg)
}

func run(cfg *configuration.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Setup signal catching for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
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
	defer func(storage *db.Storage) {
		err := storage.Close()
		if err != nil {
			fmt.Printf("error during closing storage %s", err)
		}
	}(storage)
	go func() {
		if err := server.Serve(); err != nil {
			zlogger.Log.Error("Failed to start server", zap.Error(err))
		}
	}()
	// Wait for a termination signal
	sig := <-signalChan
	fmt.Printf("Received shutdown signal %s, gracefully stopping...", sig.String())
	cancel()
	return nil
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
