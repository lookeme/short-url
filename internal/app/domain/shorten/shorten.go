package shorten

import (
	"context"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/models"
	"github.com/lookeme/short-url/internal/storage"
	"github.com/lookeme/short-url/internal/utils"
)

type URLService struct {
	shortenRepository storage.ShortenRepository
	cfg               *configuration.Config
	Log               *logger.Logger
}

func NewURLService(repository storage.ShortenRepository, log *logger.Logger, cfg *configuration.Config) URLService {
	return URLService{
		shortenRepository: repository,
		cfg:               cfg,
		Log:               log,
	}
}

func (s *URLService) CreateAndSave(originURL string, userID int) (string, error) {
	token := utils.NewShortToken(7)
	key := token.Get()
	if err := s.shortenRepository.Save(utils.CreateShortURL(key, s.cfg.Network.BaseURL), originURL, userID); err != nil {
		return "", err
	}
	return utils.CreateShortURL(key, s.cfg.Network.BaseURL), nil
}

func (s *URLService) CreateAndSaveBatch(urls []models.BatchRequest) ([]models.BatchResponse, error) {
	var dataToSave []models.ShortenData
	for _, url := range urls {
		token := utils.NewShortToken(7)
		key := token.Get()
		shorten := models.ShortenData{
			OriginalURL: url.OriginalURL,
			ShortURL:    utils.CreateShortURL(key, s.cfg.Network.BaseURL),
		}
		shorten.CorrelationID = url.CorrelationID
		dataToSave = append(dataToSave, shorten)
	}
	err := s.shortenRepository.SaveAll(dataToSave)
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

func (s *URLService) FindByURL(key string) (models.ShortenData, bool) {
	shorten, ok := s.shortenRepository.FindByURL(key)
	if !ok {
		return models.ShortenData{}, false
	}
	return shorten, true
}
func (s *URLService) FindByURLs(urls []models.BatchRequest) ([]models.ShortenData, error) {
	var keys []string
	for _, url := range urls {
		keys = append(keys, url.OriginalURL)
	}
	return s.shortenRepository.FindByURLs(keys)
}

func (s *URLService) FindByKey(key string) (models.ShortenData, bool) {
	shortURL := utils.CreateShortURL(key, s.cfg.Network.BaseURL)
	shorten, ok := s.shortenRepository.FindByKey(shortURL)
	if !ok {
		return models.ShortenData{}, false
	}
	return shorten, ok
}
func (s *URLService) FindAll() ([]models.ShortenData, error) {
	result, err := s.shortenRepository.FindAll()
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *URLService) FindAllByUserID(userID int) ([]models.ShortenData, error) {
	result, err := s.shortenRepository.FindAllByUserID(userID)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *URLService) Ping(_ context.Context) error {
	return nil
}
