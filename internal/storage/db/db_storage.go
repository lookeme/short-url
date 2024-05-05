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
func (pg *Postgres) FindByURL(key string) (string, bool) {
	query := `SELECT id, short_url, original_url FROM short WHERE original_url = @originalURL`
	args := pgx.NamedArgs{
		"originalURL": key,
	}
	var data = models.ShortenData{}
	rows := pg.connPool.QueryRow(context.Background(), query, args)
	err := rows.Scan(&data)
	if err != nil {
		pg.log.Log.Error(err.Error(), zap.String("during fetching by url", key))
		return "", false
	}
	return data.ShortURL, true
}
func (pg *Postgres) FindByKey(key string) (string, bool) {
	query := `SELECT id, short_url, original_url  FROM short WHERE short_url = @shortURL`
	args := pgx.NamedArgs{
		"shortURL": key,
	}
	rows, err := pg.connPool.Query(context.Background(), query, args)
	if err != nil {
		pg.log.Log.Error(err.Error(), zap.String("during fetching by short key", key))
		return "", false
	}
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.ShortenData])
	if err != nil {
		pg.log.Log.Error(err.Error(), zap.String("during fetching by short key", key))
		return "", false
	}
	return data.OriginalURL, true
}
func (pg *Postgres) FindAll() ([][]string, error) {
	query := `SELECT id, short_url, original_url  FROM short ORDER BY date_create DESC`
	rows, err := pg.connPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	urls, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ShortenData])
	if err != nil {
		return nil, err
	}
	result := make([][]string, len(urls))
	for i, val := range urls {
		result[i] = []string{val.ShortURL, val.OriginalURL}
	}
	return result, nil
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
