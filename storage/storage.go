package storage

// Storage 存储接口
type Storage interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
	Has(key string) bool
	Keys(prefix string) ([]string, error)
	Close() error
}
