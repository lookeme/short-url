package storage

import "github.com/lookeme/short-url/internal/models"

type ShortenRepository interface {
	Save(key, value string, userID int) error
	SaveAll(urls []models.ShortenData) error
	FindByURL(key string) (models.ShortenData, bool)
	FindByURLs(keys []string) ([]models.ShortenData, error)
	FindByKey(key string) (models.ShortenData, bool)
	FindAll() ([]models.ShortenData, error)
	FindAllByUserID(userID int) ([]models.ShortenData, error)
	Close() error
}

type UserRepository interface {
	SaveUser(name, pass string) (int, error)
	FindByID(userID int) (models.User, error)
}
