package storage

import "github.com/lookeme/short-url/internal/models"

type Repository interface {
	Save(key, value string) error
	SaveAll(urls []models.ShortenData) error
	FindByURL(key string) (models.ShortenData, bool)
	FindByURLs(keys []string) ([]models.ShortenData, error)
	FindByKey(key string) (models.ShortenData, bool)
	FindAll() ([]models.ShortenData, error)
	Close() error
}
