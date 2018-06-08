package storage

import (
	"github.com/google/uuid"
	"github.com/ashriths/go-graph/graph"
)

type InMemoryIOMapper struct {
	Memory *MemStorage
}

func (InMemoryIOMapper) StoreVertex(vertex *graph.Vertex, uid *uuid.UUID) error {
	logln(vertex)
	var u uuid.UUID
	u, e := uuid.NewUUID()
	*uid = u
	return e
}

func (InMemoryIOMapper) StoreEdge(edge *graph.Edge, uid *uuid.UUID) error {
	panic("implement me")
}

func (InMemoryIOMapper) UpdateProperties(element *graph.Element, success *bool) error {
	panic("implement me")
}

func (InMemoryIOMapper) GetElementProperties(elementId uuid.UUID, properties *interface{}) error {
	panic("implement me")
}

func (InMemoryIOMapper) RemoveVertex(vertex uuid.UUID, succ *bool) error {
	panic("implement me")
}

func (InMemoryIOMapper) RemoveEdge(edge uuid.UUID, succ *bool) error {
	panic("implement me")
}

func NewInMemoryIOMapper() *InMemoryIOMapper {
	return &InMemoryIOMapper{
		Memory: NewMemStorage(),
	}
}



var _ IOMapper = new(InMemoryIOMapper)
 
