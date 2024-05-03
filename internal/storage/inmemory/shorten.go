package inmemory

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/models"
)

type Storage struct {
	urlToKey map[string]string
	keyToURL map[string]string
	mutex    sync.RWMutex
	file     *os.File
	writer   *bufio.Writer
	log      *logger.Logger
}

func NewStorage(cfg *configuration.Storage, logger *logger.Logger) (*Storage, error) {
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Storage{
		urlToKey: make(map[string]string),
		keyToURL: make(map[string]string),
		file:     file,
		writer:   bufio.NewWriter(file),
		log:      logger,
	}, nil
}

func (s *Storage) Save(key string, value string) error {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	data := models.NewShortenData(value, key)
	s.urlToKey[value] = key
	s.keyToURL[key] = value
	if err := s.writeToFile(data); err != nil {
		return err
	}
	return nil
}

func (s *Storage) FindByURL(key string) (string, bool) {
	defer s.mutex.RUnlock()
	s.mutex.RLock()
	value, ok := s.urlToKey[key]
	return value, ok
}
func (s *Storage) FindByKey(key string) (string, bool) {
	defer s.mutex.RUnlock()
	s.mutex.RLock()
	value, ok := s.keyToURL[key]
	return value, ok
}

func (s *Storage) FindAll() ([][]string, error) {
	var result [][]string
	s.mutex.RLock()
	for k, v := range s.keyToURL {
		result = append(result, []string{k, v})
	}
	return result, nil
}

func (s *Storage) writeToFile(shortenData *models.ShortenData) error {
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

func (s *Storage) RecoverFromFile() error {
	s.log.Log.Info("Starting recovering data from file....")
	sc := bufio.NewScanner(s.file)
	for sc.Scan() {
		shorten := models.ShortenData{}
		b := sc.Bytes()
		if err := json.Unmarshal(b, &shorten); err != nil {
			return err
		}
		s.log.Log.Info("writing ...", zap.String("shorten :", shorten.String()))
		if err := s.Save(shorten.ShortURL, shorten.OriginalURL); err != nil {
			s.log.Log.Error("Error during saving ", zap.String("data", shorten.String()))
		}
	}
	return nil
}

func (s *Storage) Close() {
	s.file.Close()
}
