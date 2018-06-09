package storage

import (
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/metadata"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
)

type InMemoryIOMapper struct {
	Metadata *metadata.Metadata
	Memory   *MemoryStorage
	//Disk     *DiskStorage
}

func NewInMemoryIOMapper() *InMemoryIOMapper {
	return &InMemoryIOMapper{
		Memory: NewMemoryStorage(),
	}
}

func (self *InMemoryIOMapper) StoreVertex(vertex *graph.Vertex, success *bool) error {
	system.Logf(">> Request to Store %s\n", vertex)
	*success = false
	if e := self.Memory.StoreElement(vertex); e != nil {
		return e
	}
	if e := self.Memory.Index.CreateVertexIndex(vertex); e != nil {
		return e
	}
	*success = true
	return nil
}

func (self *InMemoryIOMapper) GetVertexById(vertexId uuid.UUID, vertex *graph.Vertex) error {
	e := self.Memory.GetVertexById(vertexId, vertex)
	system.Logln(vertex)
	return e
}

func (self *InMemoryIOMapper) GetEdgeById(edgeId uuid.UUID, vertex *graph.Edge) error {
	panic("todo")
}

func (self *InMemoryIOMapper) StoreEdge(edge *graph.Edge, success *bool) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) UpdateProperties(element *graph.Element, success *bool) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) RemoveVertex(vertex uuid.UUID, succ *bool) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) RemoveEdge(edge uuid.UUID, succ *bool) error {
	panic("implement me")
}

var _ IOMapper = new(InMemoryIOMapper)
