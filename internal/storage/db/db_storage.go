package db

import (
	"context"
	"sync"

	"github.com/lookeme/short-url/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
)

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

type Postgres struct {
	connPool *pgxpool.Pool
	log      *logger.Logger
}

type Storage struct {
	UserRepository    storage.UserRepository
	ShortenRepository storage.ShortenRepository
}

func (s Storage) Close() {
	s.Close()
}

func NewStorage(userRepo storage.UserRepository, shortRepo storage.ShortenRepository) *Storage {
	return &Storage{
		UserRepository:    userRepo,
		ShortenRepository: shortRepo,
	}
}

func (pg *Postgres) Close() error {
	pg.connPool.Close()
	return nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.connPool.Ping(ctx)
}

func New(ctx context.Context, log *logger.Logger, cfg *configuration.Storage) (*Postgres, error) {
	log.Log.Info("creating pool of conn to db...", zap.String("connString", cfg.ConnString))
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, cfg.ConnString)
		if err != nil {
			log.Log.Error(err.Error())
		}
		pgInstance = &Postgres{db, log}
	})
	err := StartMigration(pgInstance.connPool)
	if err != nil {
		return nil, err
	}
	return pgInstance, nil
}
