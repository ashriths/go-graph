package storage

import (
	"github.com/google/uuid"
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/metadata"
)

type InMemoryIOMapper struct {
	Metadata *metadata.Metadata
	Memory *MemStorage
}

func (InMemoryIOMapper) StoreVertex(vertex *graph.Vertex, uid *uuid.UUID) error {
	logf(">> Request to Store %s\n", vertex)
	panic("todo")
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
 
