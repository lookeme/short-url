package storage

type Repository interface {
	Save(key, value string) error
	FindByURL(key string) (string, bool)
	FindByKey(key string) (string, bool)
	FindAll() ([][]string, error)
	Close() error
}
