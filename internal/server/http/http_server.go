package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/server/handler"
	"github.com/lookeme/short-url/internal/storage/inmemory"
	"log"
	"net/http"
	"os"
)

var (
	ConfigFile = os.Getenv("CONFIG_FILE")
)

type Server struct {
	handler *handler.URLHandler
	config  *configuration.NetworkCfg
}

func (s *Server) Serve() error {
	r := chi.NewRouter()
	r.Post("/", s.handler.HandlePOST)
	r.Get("/{id}", s.handler.HandleGet)
	fmt.Println("Starting server on ")
	return http.ListenAndServe(":8080", r)
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
	urlService := shorten.NewURLService(storage)
	urlHandler := handler.NewURLHandler(urlService, cfg.Network)
	server := Server{
		handler: urlHandler,
		config:  cfg.Network,
	}
	return server.Serve()
}
