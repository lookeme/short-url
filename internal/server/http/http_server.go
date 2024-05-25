package http

import (
	"github.com/lookeme/short-url/internal/security"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/lookeme/short-url/internal/compression"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/server/handler"
)

type Server struct {
	handler *handler.URLHandler
	config  *configuration.NetworkCfg
	logger  *logger.Logger
	gzip    *compression.Compressor
	auth    *security.Authorization
}

func NewServer(
	handler *handler.URLHandler,
	cfg *configuration.NetworkCfg,
	logger *logger.Logger,
	compressor *compression.Compressor,
	auth *security.Authorization,
) *Server {
	return &Server{
		handler: handler,
		config:  cfg,
		logger:  logger,
		gzip:    compressor,
		auth:    auth,
	}
}

func (s *Server) Serve() error {
	r := chi.NewRouter()
	r.Use(s.logger.Middleware)
	r.Use(s.gzip.GzipMiddleware)
	r.Group(func(subRouter chi.Router) {
		subRouter.Use(s.auth.AuthMiddleware)
		subRouter.Post("/", s.handler.HandlePOST)
		subRouter.Get("/api/user/urls", s.handler.HandleUserURLs)
		subRouter.Delete("/api/user/urls", s.handler.HandleDeleteURLs)
	})
	//r.Post("/", s.handler.HandlePOST)
	r.Get("/api/user/urls", s.handler.HandleUserURLs)
	r.Post("/api/shorten", s.handler.HandleShorten)
	r.Post("/api/shorten/batch", s.handler.HandleShortenBatch)
	r.Get("/{id}", s.handler.HandleGet)
	r.Get("/ping", s.handler.HandlePing)
	r.Get("/api/user/urls", s.handler.HandleUserURLs)
	s.logger.Log.Info("shorten url service ", zap.String("starting serving on ....", s.config.ServerAddress))
	return http.ListenAndServe(s.config.ServerAddress, r)
}
