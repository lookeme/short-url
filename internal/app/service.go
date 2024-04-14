package app

type ShortenURLService interface {
	CreateAndSave(key string) (string, error)
	FindByURL(key string) (string, bool)
	FindByKey(key string) (string, bool)
}
