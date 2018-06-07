package storage

type DiskStorage struct {
}

func (DiskStorage) Clock(atLeast uint64) (error, uint64) {
	panic("implement me")
}

func (DiskStorage) Get(key string) (error, string) {
	panic("implement me")
}

func (DiskStorage) Set(key string, value string) error {
	panic("implement me")
}

func (DiskStorage) Keys(p Pattern) (error, []string) {
	panic("implement me")
}

func NewDiskStorage() *DiskStorage {
	return &DiskStorage{}
}



var _ Storage = new(DiskStorage)
