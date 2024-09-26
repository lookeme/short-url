package http

import (
	"context"
	"fmt"
	pb "github.com/lookeme/short-url/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"

	"github.com/lookeme/short-url/internal/security"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/lookeme/short-url/internal/compression"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/server/handler"
)

// Server represents a server that handles HTTP requests.
type Server struct {
	handler *handler.URLHandler
	config  *configuration.NetworkCfg
	logger  *logger.Logger
	gzip    *compression.Compressor
	auth    *security.Authorization
}

// GRPCServer represents a server that handles HTTP requests.
type GRPCServer struct {
	service pb.IShortenService
	logger  *logger.Logger
	auth    *security.Authorization
}

func NewGRPCServer(
	service pb.IShortenService,
	logger *logger.Logger,
	auth *security.Authorization) *GRPCServer {
	return &GRPCServer{
		service: service,
		logger:  logger,
		auth:    auth,
	}
}

func (g *GRPCServer) Serve() error {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	grpc.UnaryInterceptor(g.GetUnaryInterceptor())
	pb.RegisterShortenURLServiceServer(s, &g.service)
	fmt.Println("Сервер gRPC начал работу")
	return s.Serve(listen)
}

// NewServer creates a new instance of the Server struct.
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

// Serve runs the HTTP server and listens for incoming requests.
func (s *Server) Serve() error {
	r := chi.NewRouter()
	r.Use(s.logger.Middleware)
	r.Use(s.gzip.GzipMiddleware)
	r.Group(func(subRouter chi.Router) {
		subRouter.Use(s.auth.AuthMiddleware)
		subRouter.Post("/", s.handler.HandlePOST)
		subRouter.Get("/api/user/urls", s.handler.HandleUserURLs)
	})
	r.Delete("/api/user/urls", s.handler.HandleDeleteURLs)
	r.Get("/api/user/urls", s.handler.HandleUserURLs)
	r.Post("/api/shorten", s.handler.HandleShorten)
	r.Post("/api/shorten/batch", s.handler.HandleShortenBatch)
	r.Get("/{id}", s.handler.HandleGet)
	r.Get("/ping", s.handler.HandlePing)
	r.Get("/api/user/urls", s.handler.HandleUserURLs)
	s.logger.Log.Info("shorten url service ", zap.String("starting serving on ....", s.config.ServerAddress))
	return http.ListenAndServe(s.config.ServerAddress, r)
}

// Auth Interceptor
func authInterceptor(ctx context.Context, auth *security.Authorization) (context.Context, error) {
	// Extract the token from the metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	tokenArr := md["authorization"]
	if len(tokenArr) == 0 {
		usr, err := auth.UserService.CreateUser()
		if err != nil {
			return nil, fmt.Errorf("can't create token: %s", err)
		}
		token, err := auth.BuildJWTString(usr.UserID)
		metadata.AppendToOutgoingContext(ctx, "authorization", token)
	}
	token := tokenArr[0]
	if !auth.VerifyToken(token) {
		return nil, fmt.Errorf("invalid token")
	}
	return ctx, nil
}

// GetUnaryInterceptor Unary Interceptor to handle authentication
func (g *GRPCServer) GetUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Authenticate the request
		if _, err := authInterceptor(ctx, g.auth); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}
