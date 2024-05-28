package inmemory

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/models"
)

type InMemShortenStorage struct {
	urlToKey map[string]models.ShortenData
	keyToURL map[string]models.ShortenData
	id       int64
	mutex    sync.RWMutex
	file     *os.File
	writer   *bufio.Writer
	log      *logger.Logger
}

type InMemUserStorage struct {
	userMap map[int]models.User
	id      int
	mutex   sync.RWMutex
	log     *logger.Logger
}

func NewInMemUserStorage(logger *logger.Logger) (*InMemUserStorage, error) {
	return &InMemUserStorage{
		userMap: make(map[int]models.User),
		id:      0,
		log:     logger,
	}, nil
}

func (s *InMemShortenStorage) FindAllByUserID(userID int) ([]models.ShortenData, error) {
	return []models.ShortenData{}, nil
}

func NewInMemShortenStorage(cfg *configuration.Storage, logger *logger.Logger) (*InMemShortenStorage, error) {
	logger.Log.Info("Creating local storage")
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &InMemShortenStorage{
		urlToKey: make(map[string]models.ShortenData),
		keyToURL: make(map[string]models.ShortenData),
		file:     file,
		writer:   bufio.NewWriter(file),
		log:      logger,
		id:       0,
	}, nil
}

func (s *InMemShortenStorage) Save(key, value string, userID int) error {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	s.id += 1
	data := models.NewShortenData(s.id, value, key, userID)
	s.urlToKey[value] = *data
	s.keyToURL[key] = *data
	if err := s.writeToFile(data); err != nil {
		return err
	}
	return nil
}

func (s *InMemShortenStorage) SaveAll(data []models.ShortenData) error {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	for _, shorten := range data {
		s.id += 1
		shorten.ID = s.id
		s.urlToKey[shorten.OriginalURL] = shorten
		s.keyToURL[shorten.ShortURL] = shorten
		if err := s.writeToFile(&shorten); err != nil {
			return err
		}
	}
	return nil
}

func (s *InMemShortenStorage) FindByURLs(keys []string) ([]models.ShortenData, error) {
	defer s.mutex.RUnlock()
	var result []models.ShortenData
	s.mutex.RLock()
	for _, key := range keys {
		value, ok := s.urlToKey[key]
		if ok {
			result = append(result, value)
		}
	}
	return result, nil
}

func (s *InMemShortenStorage) FindByURL(key string) (models.ShortenData, bool) {
	defer s.mutex.RUnlock()
	s.mutex.RLock()
	value, ok := s.urlToKey[key]
	return value, ok
}
func (s *InMemShortenStorage) FindByKey(key string) (models.ShortenData, bool) {
	defer s.mutex.RUnlock()
	s.mutex.RLock()
	value, ok := s.keyToURL[key]
	return value, ok
}

func (s *InMemShortenStorage) FindAll() ([]models.ShortenData, error) {
	var result []models.ShortenData
	s.mutex.RLock()
	for _, shorten := range s.keyToURL {
		result = append(result, shorten)
	}
	return result, nil
}

func (s *InMemShortenStorage) writeToFile(shortenData *models.ShortenData) error {
	b, err := json.Marshal(&shortenData)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = s.writer.Write(b)
	if err != nil {
		return err
	}
	return s.writer.Flush()
}

func (s *InMemShortenStorage) RecoverFromFile() error {
	s.log.Log.Info("Starting recovering data from file....")
	sc := bufio.NewScanner(s.file)
	for sc.Scan() {
		s.id += 1
		shorten := models.ShortenData{ID: s.id}
		b := sc.Bytes()
		if err := json.Unmarshal(b, &shorten); err != nil {
			s.log.Log.Error(err.Error(), zap.String("shorten :", shorten.String()))
		}
		s.log.Log.Info("writing ...", zap.String("shorten :", shorten.String()))
		if err := s.Save(shorten.ShortURL, shorten.OriginalURL, shorten.UserID); err != nil {
			s.log.Log.Error("Error during saving ", zap.String("data", shorten.String()))
		}
	}
	return nil
}

func (s *InMemShortenStorage) Close() error {
	return s.file.Close()
}

func (s *InMemUserStorage) SaveUser(name, pass string) (int, error) {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	s.id += 1
	user := models.User{
		UserID: s.id,
		Name:   name,
		Pass:   pass,
	}
	s.userMap[s.id] = user
	return user.UserID, nil
}

func (s *InMemUserStorage) FindByID(userID int) (models.User, error) {
	user, ok := s.userMap[userID]
	if !ok {
		return user, errors.New("user doesn't exist")
	}
	return user, nil
}

func (s *InMemShortenStorage) DeleteByShortURL(shortURL string) bool {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	val := s.keyToURL[shortURL]
	val.DeletedFlag = true
	s.keyToURL[shortURL] = val
	return true
}
