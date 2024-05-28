package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/lookeme/short-url/internal/models"
	"go.uber.org/zap"
)

type ShortenRepository struct {
	postgres *Postgres
}

func NewShortenRepository(postgres *Postgres) *ShortenRepository {
	return &ShortenRepository{
		postgres: postgres,
	}
}

func (r *ShortenRepository) Save(key, value string, userID int) error {
	query := `INSERT INTO short (original_url, short_url, user_id) VALUES (@originalURL, @shortURL, @userID)`
	args := pgx.NamedArgs{
		"originalURL": value,
		"shortURL":    key,
		"userID":      userID,
	}
	_, err := r.postgres.connPool.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}
	return nil
}
func (r *ShortenRepository) FindByURL(key string) (models.ShortenData, bool) {
	query := `SELECT id, correlation_id, short_url, original_url, user_id FROM short WHERE original_url = @originalURL`
	args := pgx.NamedArgs{
		"originalURL": key,
	}
	var data models.ShortenData
	row, err := r.postgres.connPool.Query(context.Background(), query, args)
	if err != nil {
		return data, false
	}
	data, err = pgx.CollectOneRow(row, pgx.RowToStructByPos[models.ShortenData])
	if err != nil {
		r.postgres.log.Log.Error(err.Error(), zap.String("during fetching by url", key))
		return data, false
	}
	return data, true
}

func (r *ShortenRepository) FindByURLs(keys []string) ([]models.ShortenData, error) {
	query := `SELECT id, short_url, original_url, correlation_id, user_id FROM short WHERE original_url = ANY (@originalURL)`
	args := pgx.NamedArgs{
		"originalURL": keys,
	}
	rows, err := r.postgres.connPool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ShortenData])
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (r *ShortenRepository) FindByKey(key string) (models.ShortenData, bool) {
	query := `SELECT id, correlation_id, short_url, original_url, user_id  FROM short WHERE short_url = @shortURL`
	args := pgx.NamedArgs{
		"shortURL": key,
	}
	row, err := r.postgres.connPool.Query(context.Background(), query, args)
	if err != nil {
		r.postgres.log.Log.Error(err.Error(), zap.String("during fetching by short key", key))
		return models.ShortenData{}, false
	}
	data, err := pgx.CollectOneRow(row, pgx.RowToStructByPos[models.ShortenData])
	if err != nil {
		r.postgres.log.Log.Error(err.Error(), zap.String("during fetching by short key", key))
		return models.ShortenData{}, false
	}
	return data, true
}
func (r *ShortenRepository) FindAll() ([]models.ShortenData, error) {
	query := `SELECT id, short_url, original_url, correlation_id, user_id FROM short ORDER BY date_create DESC`
	rows, err := r.postgres.connPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ShortenData])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ShortenRepository) SaveAll(rows []models.ShortenData) error {
	if len(rows) == 0 {
		return nil
	}
	conn, err := r.postgres.connPool.Acquire(context.Background())
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

func (r *ShortenRepository) FindAllByUserID(userID int) ([]models.ShortenData, error) {
	query := `SELECT id, short_url, original_url, correlation_id, user_id FROM short WHERE user_id = (@userID) ORDER BY date_create DESC`
	args := pgx.NamedArgs{
		"userID": userID,
	}
	rows, err := r.postgres.connPool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ShortenData])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ShortenRepository) Close() error {
	r.postgres.connPool.Close()
	return nil
}
