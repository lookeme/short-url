package main

import (
	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/server/handler"
	"github.com/lookeme/short-url/internal/server/http"
	"github.com/lookeme/short-url/internal/storage/inmemory"
	"log"
)

func main() {
	parseFlags()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	storage := inmemory.NewStorage()
	urlService := shorten.NewURLService(storage)
	urlHandler := handler.NewURLHandler(urlService, &cfg)
	server := http.Server{
		Handler: *urlHandler,
		Config:  *cfg.Network,
	}
	return server.Serve()
}
