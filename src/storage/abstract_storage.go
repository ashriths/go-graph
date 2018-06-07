package storage

type Storage interface {
	Clock(atLeast uint64) (error, uint64)
	Get(key string) (error, string)
	Set(key string, value string) error
	Keys(p Pattern) (error, []string)
}

type StorageConfig struct {
	Addr string
	Ready chan<- bool
	Store IOMapper
}

