package app

type Service struct {
	ShortenUrlService ShortenUrlService
}

type ShortenUrlService interface {
	CreateAndSave(key string) (string, error)
	FindByURL(key string) (string, bool)
	FindByKey(key string) (string, bool)
}
