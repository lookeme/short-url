package shorten

import (
	"fmt"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/models"
	"github.com/lookeme/short-url/internal/storage"
	"github.com/lookeme/short-url/internal/utils"
)

type URLService struct {
	shortenRepository storage.Repository
	cfg               *configuration.Config
}

func NewURLService(repository storage.Repository, cfg *configuration.Config) *URLService {
	return &URLService{
		shortenRepository: repository,
		cfg:               cfg,
	}
}

func (s *URLService) CreateAndSave(originURL string) (string, error) {
	key, ok := s.FindByURL(originURL)
	if !ok {
		token := utils.NewShortToken(7)
		key = token.Get()
		if err := s.shortenRepository.Save(key, originURL); err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/%s", s.cfg.Network.BaseURL, key), nil
	} else {
		return fmt.Sprintf("%s/%s", s.cfg.Network.BaseURL, key), nil
	}
}

func (s *URLService) FindByURL(key string) (string, bool) {
	return s.shortenRepository.FindByURL(key)
}

func (s *URLService) FindByKey(key string) (string, bool) {
	return s.shortenRepository.FindByKey(key)
}
func (s *URLService) FindAll() ([]models.ShortenData, error) {
	var result []models.ShortenData
	data, err := s.shortenRepository.FindAll()
	if err != nil {
		return result, err
	}
	for _, v := range data {
		result = append(result, models.ShortenData{
			ShortURL:    fmt.Sprintf("%s/%s", s.cfg.Network.BaseURL, v[0]),
			OriginalURL: v[1],
		})
	}
	return result, nil
}
