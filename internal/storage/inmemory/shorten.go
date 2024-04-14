package inmemory

import "sync"

type Storage struct {
	urlToKey  map[string]string
	keyTouUrl map[string]string
	mutex     sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		urlToKey:  make(map[string]string),
		keyTouUrl: make(map[string]string),
	}
}

func (s *Storage) Save(key, value string) error {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	s.urlToKey[key] = value
	s.keyTouUrl[value] = key
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
	value, ok := s.keyTouUrl[key]
	return value, ok
}
