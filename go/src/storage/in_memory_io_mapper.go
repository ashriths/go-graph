package storage

import "github.com/google/uuid"

type InMemoryIOMapper struct {
	Memory *MemStorage
}

func NewInMemoryIOMapper() *InMemoryIOMapper {
	return &InMemoryIOMapper{
		Memory: NewMemStorage(),
	}
}

func (self *InMemoryIOMapper) StoreVertex(properties interface{}) (error, uuid.UUID) {
	panic("implement me")
}

func (self *InMemoryIOMapper) StoreEdge(srcVertex uuid.UUID, destVertex uuid.UUID, properties interface{}) (error, uuid.UUID) {
	panic("implement me")
}

func (self *InMemoryIOMapper) UpdateProperties(elementId uuid.UUID, properties interface{}) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) UpdatePropertyByName(elementId uuid.UUID, key string, value string) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) GetElementProperties(elementId uuid.UUID) (error, interface{}) {
	panic("implement me")
}

func (self *InMemoryIOMapper) GetElementPropertyByName(elementId uuid.UUID, key string) (error, interface{}) {
	panic("implement me")
}

func (self *InMemoryIOMapper) RemoveVertex(vertex uuid.UUID) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) RemoveEdge(edge uuid.UUID) error {
	panic("implement me")
}

var _ IOMapper = new(InMemoryIOMapper)
 
