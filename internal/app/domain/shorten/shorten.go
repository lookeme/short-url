package shorten

import (
	"github.com/lookeme/short-url/internal/storage"
	"github.com/lookeme/short-url/internal/utils"
)

type UrlService struct {
	shortenRepository storage.Repository
}

func NewUrlService(repository storage.Repository) *UrlService {
	return &UrlService{
		shortenRepository: repository,
	}
}

func (s *UrlService) CreateAndSave(key string) (string, error) {
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

func (s *UrlService) FindByURL(key string) (string, bool) {
	return s.shortenRepository.FindByURL(key)
}

func (s *UrlService) FindByKey(key string) (string, bool) {
	return s.shortenRepository.FindByKey(key)
}
