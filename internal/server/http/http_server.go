package http

import (
	"fmt"
	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/configuration"
	shorten2 "github.com/lookeme/short-url/internal/server/http/handler/shorten"
	"github.com/lookeme/short-url/internal/storage/inmemory"
	"log"
	"net/http"
	"os"
)

var (
	ConfigFile = os.Getenv("CONFIG_FILE")
)

type Server struct {
	handler *shorten2.UrlHandler
	config  *configuration.NetworkCfg
}

func (s *Server) Serve() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, s.handler.Index)
	var host string
	if s.config == nil {
		host = ":8080"
	}
	fmt.Println("Starting server on ", host)
	return http.ListenAndServe(host, mux)
}

func Run() error {
	cfg := &configuration.Config{}
	if ConfigFile != "" {
		var err error
		cfg, err = configuration.LoadCfg(ConfigFile)
		if err != nil {
			log.Fatal("Can't find configuration file")
		}
	}
	storage := inmemory.NewStorage()
	urlService := shorten.NewUrlService(storage)
	urlHandler := shorten2.NewUrlHandler(urlService, cfg.Network)
	server := Server{
		handler: urlHandler,
		config:  cfg.Network,
	}
	return server.Serve()
}
