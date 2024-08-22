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

	// pgInstance is a pointer to a Postgres struct variable used for storing the initialized instance of the Postgres database connection pool and logger.
	pgInstance *Postgres

	// pgOnce is a sync.Once variable used for lazy initialization of the pgInstance variable in the New function.
	pgOnce sync.Once
)

// Postgres is a type representing a connection pool to a PostgreSQL database.
// It contains a *pgxpool.Pool object for managing connections and a *logger.Logger object for logging.
type Postgres struct {
	connPool *pgxpool.Pool
	log      *logger.Logger
}

// Implementation omitted for brevity
type Storage struct {
	UserRepository    storage.UserRepository
	ShortenRepository storage.ShortenRepository
}

// Close closes the storage by calling the Close method of the underlying ShortenRepository.
// It returns an error if there is an error closing the repository.
func (s *Storage) Close() error {
	return s.ShortenRepository.Close()
}

// NewStorage creates a new instance of Storage by accepting an implementation of UserRepository and ShortenRepository.
// It initializes the UserRepository and ShortenRepository fields of the Storage struct with the provided instances.
// It returns a pointer to the newly created Storage instance.
// Example usage:
// userRepo, err := storage.NewInMemUserStorage(log)
//
//	if err != nil {
//	    log.Error("Failed to create user repository", zap.Error(err))
//	    return nil, err
//	}
//
// shortenRepo, err := storage.NewInMemShortenStorage(cfg, log)
//
//	if err != nil {
//	    log.Error("Failed to create shorten repository", zap.Error(err))
//	    return nil, err
//	}
//
// storage := NewStorage(userRepo, shortenRepo)
func NewStorage(userRepo storage.UserRepository, shortRepo storage.ShortenRepository) *Storage {
	return &Storage{
		UserRepository:    userRepo,
		ShortenRepository: shortRepo,
	}
}

// Close closes the PostgreSQL connection pool by calling the Close method of the connPool.
// It does not return any error.
func (pg *Postgres) Close() error {
	pg.connPool.Close()
	return nil
}

// Ping pings the Postgres database by calling the Ping method of the connPool.
// It takes a context.Context as an argument and returns an error if there is an error pinging the database.
// Example usage:
//
//	err := pg.Ping(ctx)
func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.connPool.Ping(ctx)
}

// New creates a new instance of Postgres by accepting a context, logger, and storage configuration.
// It logs an informational message indicating the creation of the connection pool.
// It uses sync.Once to ensure that the connection pool is created only once.
// Inside the sync.Once's Do() function, it creates a new pgxpool.Pool with the provided context and connection string.
// If there is an error during the creation of the connection pool, it logs the error and returns it.
// After the connection pool is created, it assigns it to the pgInstance variable.
// It then calls the StartMigration function to initialize the database schema using goose.
// If there is an error during the migration, it returns the error.
// Finally, it returns the pgInstance and nil error.
// Example usage:
// ctx := context.Background()
// log := &logger.Logger{Log: zap.NewNop()}
// cfg := &configuration.Storage{ConnString: "postgres://localhost:5432/mydb"}
// pg, err := New(ctx, log, cfg)
//
//	if err != nil {
//	    log.Log.Error("Failed to create Postgres instance", zap.Error(err))
//	    return nil, err
//	}
//
// defer pg.Close()
// err = pg.Ping(ctx)
//
//	if err != nil {
//	    log.Log.Error("Failed to ping database", zap.Error(err))
//	    return nil, err
//	}
//
// // Perform database operations using the Postgres instance
// ...
// // Close the database connection pool when done
// err = pg.Close()
//
//	if err != nil {
//	    log.Log.Error("Failed to close database connection pool", zap.Error(err))
//	    return err
//	}
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
