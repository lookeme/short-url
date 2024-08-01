package shorten

import (
	"context"
	"sync"

	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/models"
	"github.com/lookeme/short-url/internal/storage"
	"github.com/lookeme/short-url/internal/utils"
	"go.uber.org/zap"
)

// URLService is a type that provides
type URLService struct {
	shortenRepository storage.ShortenRepository
	cfg               *configuration.Config
	Log               *logger.Logger
}

// NewURLService creates a new instance of URLService by initializing the shorten repository, configuration, and logger.
// It returns the created URLService.
func NewURLService(repository storage.ShortenRepository, log *logger.Logger, cfg *configuration.Config) URLService {
	return URLService{
		shortenRepository: repository,
		cfg:               cfg,
		Log:               log,
	}
}

// handle error
func (s *URLService) CreateAndSave(originURL string, userID int) (string, error) {
	token := utils.NewShortToken(7)
	key := token.Get()
	if err := s.shortenRepository.Save(utils.CreateShortURL(key, s.cfg.Network.BaseURL), originURL, userID); err != nil {
		return "", err
	}
	return utils.CreateShortURL(key, s.cfg.Network.BaseURL), nil
}

// CreateAndSaveBatch takes a slice of BatchRequest and creates and saves a batch of ShortenData objects.
// It generates a short token for each URL and creates a ShortenData object with the original URL and the short URL.
// The correlation ID from each BatchRequest is copied to the corresponding ShortenData object.
// The generated ShortenData objects are then saved using the shortenRepository's SaveAll method.
// Finally, it creates a slice of BatchResponse objects with the correlation ID and short URL from each
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

// FindByURL searches for a ShortenData object in the shortenRepository based on the given key.
// If the object is found, it returns the ShortenData object and true. Otherwise, it returns an empty ShortenData object and false.
func (s *URLService) FindByURL(key string) (models.ShortenData, bool) {
	shorten, ok := s.shortenRepository.FindByURL(key)
	if !ok {
		return models.ShortenData{}, false
	}
	return shorten, true
}

// FindByURLs retrieves a batch of ShortenData objects from the shortenRepository based on the specified URLs.
// It extracts the original URLs from the batch request and uses them as keys to query the repository.
// The method returns the matching ShortenData objects and any error that occurred during the query.
func (s *URLService) FindByURLs(urls []models.BatchRequest) ([]models.ShortenData, error) {
	var keys []string
	for _, url := range urls {
		keys = append(keys, url.OriginalURL)
	}
	return s.shortenRepository.FindByURLs(keys)
}

// FindByKey finds the shorten data with the given key in the URLService.
// It creates the short URL using the key and the base URL from the configuration.
func (s *URLService) FindByKey(key string) (models.ShortenData, bool) {
	shortURL := utils.CreateShortURL(key, s.cfg.Network.BaseURL)
	shorten, ok := s.shortenRepository.FindByKey(shortURL)
	if !ok {
		return models.ShortenData{}, false
	}
	return shorten, ok
}

// FindAll retrieves all shorten data from the repository.
// It returns a slice of ShortenData and an error if any.
func (s *URLService) FindAll() ([]models.ShortenData, error) {
	result, err := s.shortenRepository.FindAll()
	if err != nil {
		return result, err
	}
	return result, nil
}

// FindAllByUserID retrieves all shorten data associated with a specific user by their userID.
// It returns a slice of models.ShortenData and an error.
func (s *URLService) FindAllByUserID(userID int) ([]models.ShortenData, error) {
	result, err := s.shortenRepository.FindAllByUserID(userID)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Ping is a method of the URLService struct that is used to ping the service and check if it is available.
// It takes a context.Context object as a parameter, but it is not used in the implementation.
// It returns an error if there is an error during
func (s *URLService) Ping(_ context.Context) error {
	return nil
}

// DeleteByShortURLs deletes URLs based on the provided shortURLs.
func (s *URLService) DeleteByShortURLs(shortURLs []string) error {
	results := make(chan bool)
	var wg sync.WaitGroup

	for _, val := range shortURLs {
		wg.Add(1)
		url := val
		go func() {
			defer wg.Done()
			results <- s.shortenRepository.DeleteByShortURL(utils.CreateShortURL(url, s.cfg.Network.BaseURL))
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	for i := range results {
		s.Log.Log.Info("delete operation", zap.Bool("val", i))
	}
	return nil
}
