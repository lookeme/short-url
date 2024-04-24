package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/server/handler"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	Handler *handler.URLHandler
	Config  *configuration.NetworkCfg
	Logger  *logger.Logger
}

func (s *Server) Serve() error {
	r := chi.NewRouter()
	r.Use(s.Logger.Middleware)
	r.Post("/", s.Handler.HandlePOST)
	r.Post("/api/shorten", s.Handler.HandleShorten)
	r.Get("/{id}", s.Handler.HandleGet)
	s.Logger.Log.Info("shorten url service ", zap.String("starting serving on ....", s.Config.ServerAddress))
	return http.ListenAndServe(s.Config.ServerAddress, r)
}
