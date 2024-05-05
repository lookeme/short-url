package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
)

type Storage struct {
	connPool *pgxpool.Pool
	log      *logger.Logger
	connStr  string
}

func (s *Storage) Close() {
	s.connPool.Close()
}

func (s *Storage) Save(key, value string) error {
	return nil
}
func (s *Storage) FindByURL(key string) (string, bool) {
	return "", false
}
func (s *Storage) FindByKey(key string) (string, bool) {
	return "", false
}
func (s *Storage) FindAll() ([][]string, error) {
	return nil, nil
}

func (s *Storage) Ping(ctx context.Context) error {
	//connection, err := s.connPool.Acquire(ctx)
	//defer connection.Release()
	//if err != nil {
	//	return err
	//}
	//return s.connPool.Ping(ctx)
	conn, err := pgx.Connect(context.Background(), s.connStr)

	if err != nil {
		s.log.Log.Error(err.Error())
		return err
	}
	defer conn.Close(ctx)
	return conn.Ping(ctx)
}

func NewDBStorage(ctx context.Context, log *logger.Logger, cfg *configuration.Storage) (*Storage, error) {
	log.Log.Info("creating pool of conn to db...", zap.String("connString", cfg.ConnString))
	connPool, err := pgxpool.New(ctx, cfg.ConnString)
	if err != nil {
		return nil, err
	}
	return &Storage{connPool: connPool, log: log, connStr: cfg.ConnString}, nil
}
