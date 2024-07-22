// Package storage declares interfaces for data access methods related to URL shortening and user operations.
package storage

import "github.com/lookeme/short-url/internal/models"

// ShortenRepository interface represents the necessary CRUD operations for handling URLs in persistence storage.
type ShortenRepository interface {
	Save(key, value string, userID int) error
	SaveAll(urls []models.ShortenData) error
	FindByURL(key string) (models.ShortenData, bool)
	FindByURLs(keys []string) ([]models.ShortenData, error)
	FindByKey(key string) (models.ShortenData, bool)
	FindAll() ([]models.ShortenData, error)
	FindAllByUserID(userID int) ([]models.ShortenData, error)
	Close() error
	DeleteByShortURL(shortURL string) bool
}

// UserRepository interface defines the methods necessary for handling users in persistence storage.
type UserRepository interface {
	SaveUser(name, pass string) (int, error)
	FindByID(userID int) (models.User, error)
}
