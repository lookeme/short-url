package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/server/handler"
	"net/http"
	"os"
)

var (
	ConfigFile = os.Getenv("CONFIG_FILE")
)

type Server struct {
	Handler handler.URLHandler
	Config  configuration.NetworkCfg
}

func (s *Server) Serve() error {
	r := chi.NewRouter()
	r.Post("/", s.Handler.HandlePOST)
	r.Get("/{id}", s.Handler.HandleGet)
	fmt.Printf("Starting server on port %s", s.Config.ServerAddress)
	return http.ListenAndServe(s.Config.ServerAddress, r)
}
