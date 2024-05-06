package shorten

import (
	"context"
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

func NewURLService(repository storage.Repository, cfg *configuration.Config) URLService {
	return URLService{
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

func (s *URLService) CreateAndSaveBatch(urls []models.BatchRequest) ([]models.BatchResponse, error) {
	shortenData, err := s.FindByURLs(urls)
	if err != nil {
		return nil, err
	}
	var dataToSave []models.ShortenData
	for _, url := range urls {
		if !utils.Contains(shortenData, url.OriginalURL) {
			token := utils.NewShortToken(7)
			key := token.Get()
			shorten := models.ShortenData{OriginalURL: url.OriginalURL, ShortURL: key}
			shorten.CorrelationID = url.CorrelationID
			dataToSave = append(dataToSave, shorten)
		}
	}
	err = s.shortenRepository.SaveAll(dataToSave)
	if err != nil {
		return nil, err
	}
	var result []models.BatchResponse
	for _, shorten := range dataToSave {
		r := models.BatchResponse{
			CorrelationID: shorten.CorrelationID,
			ShortURL:      shorten.ShortURL,
		}
		result = append(result, r)
	}
	return result, nil
}

func (s *URLService) FindByURL(key string) (string, bool) {
	shorten, ok := s.shortenRepository.FindByURL(key)
	if !ok {
		return "", false
	}
	return shorten.ShortURL, true
}
func (s *URLService) FindByURLs(urls []models.BatchRequest) ([]models.ShortenData, error) {
	var keys []string
	for _, url := range urls {
		keys = append(keys, url.OriginalURL)
	}
	return s.shortenRepository.FindByURLs(keys)
}

func (s *URLService) FindByKey(key string) (string, bool) {
	shorten, ok := s.shortenRepository.FindByKey(key)
	if !ok {
		return "", false
	}
	return shorten.OriginalURL, ok
}
func (s *URLService) FindAll() ([]models.ShortenData, error) {
	result, err := s.shortenRepository.FindAll()
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *URLService) Ping(_ context.Context) error {
	return nil
}
