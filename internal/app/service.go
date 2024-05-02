package app

import "github.com/lookeme/short-url/internal/models"

type ShortenURLService interface {
	CreateAndSave(key string) (string, error)
	FindByURL(key string) (string, bool)
	FindByKey(key string) (string, bool)
	FindAll() ([]models.ShortenData, error)
}
