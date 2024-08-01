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

// InMemShortenStorage is an in-memory implementation of a storage for shortened URLs.
type InMemShortenStorage struct {
	urlToKey map[string]models.ShortenData
	keyToURL map[string]models.ShortenData
	id       int64
	mutex    sync.RWMutex
	file     *os.File
	writer   *bufio.Writer
	log      *logger.Logger
}

// InMemUserStorage is an in-memory implementation of a storage for user data.
// The userMap field is a map that stores users by their ID.
// The key is the user's ID (integer), and the value is a User object.
type InMemUserStorage struct {
	userMap map[int]models.User
	id      int
	mutex   sync.RWMutex
	log     *logger.Logger
}

// NewInMemUserStorage creates a new instance of InMemUserStorage with the
func NewInMemUserStorage(logger *logger.Logger) (*InMemUserStorage, error) {
	return &InMemUserStorage{
		userMap: make(map[int]models.User),
		id:      0,
		log:     logger,
	}, nil
}

// FindAllByUserID retrieves all ShortenData objects associated with a given userID.
// It returns an empty slice of ShortenData objects and a nil error.
func (s *InMemShortenStorage) FindAllByUserID(userID int) ([]models.ShortenData, error) {
	return []models.ShortenData{}, nil
}

// NewInMemShortenStorage creates a new instance of InMemShortenStorage with the given configuration and logger.
// It opens the file specified in the configuration and initializes the necessary variables.
// It returns a pointer to the created InMemShortenStorage and an error if any.
// The InMemShortenStorage struct provides various methods to interact with the shorten data storage.
// Example usage:
//
//	cfg := &configuration.Storage{
//	   FileStoragePath: "/path/to/file",
//	   ConnString:      "postgres://user:password@localhost:5432/db",
//	   PGPoolCfg:       &pgxpool
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

// Save method saves a new ShortenData object to the in-memory storage, as well as writes it to a file.
// It takes the key, value, and userID as parameters.
// The key is the shortened URL, the value is the original URL, and the userID is the ID of the user who created the shorten URL.
// It first acquires
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

// FindByURLs retrieves a slice of ShortenData objects associated with the given URLs.
// It searches the urlToKey map for each URL in the provided `keys
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

// FindByURL retrieves the ShortenData object associated with the given URL key.
// It returns the ShortenData object and a boolean value indicating whether the key was found.
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

// FindAll retrieves all ShortenData objects from the InMemShortenStorage.
// It iterates over the keyToURL map to collect all shorten data and returns them.
// It returns the result slice of ShortenData objects and a nil error.
func (s *InMemShortenStorage) FindAll() ([]models.ShortenData, error) {
	var result []models.ShortenData
	s.mutex.RLock()
	for _, shorten := range s.keyToURL {
		result = append(result, shorten)
	}
	return result, nil
}

// writeToFile writes the given ShortenData to a file.
// It first marshals the shortenData object into JSON format in
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

// RecoverFromFile reads data from a file and recovers it into the InMemShortenStorage object.
// It starts by logging an info message indicating that the recovery process has begun.
// It then creates a new scanner to read the file line by line.
// For each line, it increments the ID of the storage object and creates a new ShortenData object with that ID.
// It reads the line
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

// Close closes the file associated with the InMemShortenStorage object.
// It returns an error if the file fails to close.
func (s *InMemShortenStorage) Close() error {
	return s.file.Close()
}

// SaveUser saves a user with the given name and password into the user storage.
// It generates a unique UserID for the user and adds it to the userMap.
// The method returns the generated UserID and a nil error.
// It uses
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

// FindByID retrieves a User object based on the provided userID.
// It returns the User object and a nil error if the user
func (s *InMemUserStorage) FindByID(userID int) (models.User, error) {
	user, ok := s.userMap[userID]
	if !ok {
		return user, errors.New("user doesn't exist")
	}
	return user, nil
}

// DeleteByShortURL deletes a ShortenData object with the specified shortURL.
// It sets the DeletedFlag to true for the specified shortURL in the keyToURL map.
// It returns true to indicate that the deletion was successful.
func (s *InMemShortenStorage) DeleteByShortURL(shortURL string) bool {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	val := s.keyToURL[shortURL]
	val.DeletedFlag = true
	s.keyToURL[shortURL] = val
	return true
}
