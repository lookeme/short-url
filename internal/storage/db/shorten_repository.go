package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lookeme/short-url/internal/models"
	"go.uber.org/zap"
)

// ShortenRepository represents a repository for storing
type ShortenRepository struct {
	postgres *Postgres
}

// NewShortenRepository initializes a new instance of
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

// FindByURL searches for a record in the "short" table based on the original URL.
// It returns the matching record and a flag indicating whether the record was found.
func (r *ShortenRepository) FindByURL(key string) (models.ShortenData, bool) {
	query := `SELECT id, correlation_id, short_url, original_url, user_id, is_deleted FROM short WHERE original_url = @originalURL AND is_deleted = false`
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

// FindByURLs retrieves a list of ShortenData objects from the database by matching the original URLs with the given keys.
// It executes a SELECT query on the 'short' table.
// It returns an
func (r *ShortenRepository) FindByURLs(keys []string) ([]models.ShortenData, error) {
	query := `SELECT id, short_url, original_url, correlation_id, user_id, is_deleted FROM short WHERE original_url = ANY (@originalURL) AND is_deleted = false`
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

// FindByKey searches for a ShortenData object in the database based on a given short URL key.
// It returns the found ShortenData object and a boolean value indicating whether the data
func (r *ShortenRepository) FindByKey(key string) (models.ShortenData, bool) {
	query := `SELECT id, correlation_id, short_url, original_url, user_id, is_deleted FROM short WHERE short_url = @shortURL`
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

// FindAll retrieves all shorten data from the database that are not deleted, ordered by date created in descending order.
// It returns a slice of models.ShortenData and an error. If there is an error executing the query, the error is returned.
// If there are no results, an empty slice is returned.
// Example usage:
// shortens, err := shortenRepo.FindAll()
//
//	if err != nil {
//	    // handle error
//	}
//
//	for _, shorten := range shortens {
//	    fmt.Println(shorten)
func (r *ShortenRepository) FindAll() ([]models.ShortenData, error) {
	query := `SELECT id, short_url, original_url, correlation_id, user_id, is_deleted FROM short WHERE is_deleted = false ORDER BY date_create DESC`
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

// SaveAll saves multiple rows of ShortenData to the 'short' table in the database.
// If the rows parameter is empty, the method returns nil immediately.
// First, it acquires a connection from the connection pool.
// Then, it uses the Acquired connection to execute the CopyFrom method, which performs a bulk insert operation.
// The CopyFrom method copies the provided rows into the 'short' table.
// The CopyFromSlice method is used as a callback to convert each ShortenData struct into the required format for the CopyFrom
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

// FindAllByUserID retrieves all shorten data for a given userID that have not been deleted.
// It returns a slice of models.ShortenData and an error, if any.
func (r *ShortenRepository) FindAllByUserID(userID int) ([]models.ShortenData, error) {
	query := `SELECT id, short_url, original_url, correlation_id, user_id, is_deleted FROM short WHERE user_id = (@userID) AND short.is_deleted = false ORDER BY date_create DESC`
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

// DeleteByShortURL deletes a record from the "short" table in the database based on the short URL.
func (r *ShortenRepository) DeleteByShortURL(shortURL string) bool {
	var err error
	sqlStatement := `UPDATE short SET is_deleted = true WHERE short_url = $1`
	_, err = r.postgres.connPool.Exec(context.Background(), sqlStatement, shortURL)
	return err == nil
}

// Close closes the connection pool of the ShortenRepository's Postgres instance.
func (r *ShortenRepository) Close() error {
	r.postgres.connPool.Close()
	return nil
}
