package app

import "github.com/lookeme/short-url/internal/models"

type ShortenURLService interface {
	CreateAndSave(key string) (string, error)
	FindByURL(key string) (*models.ShortenData, bool)
	FindByKey(key string) (*models.ShortenData, bool)
	FindAll() ([]models.ShortenData, error)
	CreateAndSaveBatch(urls []models.BatchRequest) ([]models.BatchResponse, error)
}
