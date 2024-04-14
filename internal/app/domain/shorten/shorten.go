package shorten

import (
	"github.com/lookeme/short-url/internal/storage"
	"github.com/lookeme/short-url/internal/utils"
)

type URLService struct {
	shortenRepository storage.Repository
}

func NewURLService(repository storage.Repository) *URLService {
	return &URLService{
		shortenRepository: repository,
	}
}

func (s *URLService) CreateAndSave(key string) (string, error) {
	val, ok := s.FindByURL(key)
	if !ok {
		token := utils.NewShortToken(7)
		val := token.Get()
		if err := s.shortenRepository.Save(key, val); err != nil {
			return "", err
		}
		return val, nil
	} else {
		return val, nil
	}
}

func (s *URLService) FindByURL(key string) (string, bool) {
	return s.shortenRepository.FindByURL(key)
}

func (s *URLService) FindByKey(key string) (string, bool) {
	return s.shortenRepository.FindByKey(key)
}
