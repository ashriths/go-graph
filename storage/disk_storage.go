package storage

import (
	"github.com/ashriths/go-graph/graph"
)

type DiskStorage struct {
}

func (self *DiskStorage) StoreElement(element graph.ElementInterface) error {
	panic("implement me")
}

func NewDiskStorage() *DiskStorage {
	return &DiskStorage{}
}

var _ Storage = new(DiskStorage)
