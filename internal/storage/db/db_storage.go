package db

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/models"
)

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

type Postgres struct {
	connPool *pgxpool.Pool
	log      *logger.Logger
}

func (pg *Postgres) Close() error {
	pg.connPool.Close()
	return nil
}

func (pg *Postgres) Save(key, value string) error {
	query := `INSERT INTO short (original_url, short_url) VALUES (@originalURL, @shortURL)`
	args := pgx.NamedArgs{
		"originalURL": value,
		"shortURL":    key,
	}
	_, err := pg.connPool.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}
	return nil
}
func (pg *Postgres) FindByURL(key string) (models.ShortenData, bool) {
	query := `SELECT id, correlation_id, short_url, original_url FROM short WHERE original_url = @originalURL`
	args := pgx.NamedArgs{
		"originalURL": key,
	}
	var data models.ShortenData
	row, err := pg.connPool.Query(context.Background(), query, args)
	if err != nil {
		return data, false
	}
	data, err = pgx.CollectOneRow(row, pgx.RowToStructByPos[models.ShortenData])
	if err != nil {
		pg.log.Log.Error(err.Error(), zap.String("during fetching by url", key))
		return data, false
	}
	return data, true
}

func (pg *Postgres) FindByURLs(keys []string) ([]models.ShortenData, error) {
	query := `SELECT id, short_url, original_url, correlation_id FROM short WHERE original_url = ANY (@originalURL)`
	args := pgx.NamedArgs{
		"originalURL": keys,
	}
	rows, err := pg.connPool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ShortenData])
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (pg *Postgres) FindByKey(key string) (models.ShortenData, bool) {
	query := `SELECT id, correlation_id, short_url, original_url  FROM short WHERE short_url = @shortURL`
	args := pgx.NamedArgs{
		"shortURL": key,
	}
	row, err := pg.connPool.Query(context.Background(), query, args)
	if err != nil {
		pg.log.Log.Error(err.Error(), zap.String("during fetching by short key", key))
		return models.ShortenData{}, false
	}
	data, err := pgx.CollectOneRow(row, pgx.RowToStructByPos[models.ShortenData])
	if err != nil {
		pg.log.Log.Error(err.Error(), zap.String("during fetching by short key", key))
		return models.ShortenData{}, false
	}
	return data, true
}
func (pg *Postgres) FindAll() ([]models.ShortenData, error) {
	query := `SELECT id, short_url, original_url, correlation_id  FROM short ORDER BY date_create DESC`
	rows, err := pg.connPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ShortenData])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (pg *Postgres) SaveAll(rows []models.ShortenData) error {
	if len(rows) == 0 {
		return nil
	}
	conn, err := pg.connPool.Acquire(context.Background())
	defer conn.Release()
	if err != nil {
		return err
	}
	_, err = conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"short"},
		[]string{"correlation_id", "short_url", "original_url"},
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
			return []any{rows[i].CorrelationID, rows[i].ShortURL, rows[i].OriginalURL}, nil
		}),
	)
	return err
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
