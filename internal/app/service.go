package app

import "github.com/lookeme/short-url/internal/models"

type ShortenURLService interface {
	CreateAndSave(key string, userID int) (string, error)
	FindByURL(key string) (models.ShortenData, bool)
	FindByKey(key string) (models.ShortenData, bool)
	FindAll() ([]models.ShortenData, error)
	CreateAndSaveBatch(urls []models.BatchRequest) ([]models.BatchResponse, error)
	DeleteByShortURLAndUserID(urls []string, userID int) error
}

type UserService interface {
	CreateUser(userName string) (models.User, error)
	FindByID(userID int) (models.User, bool)
}
