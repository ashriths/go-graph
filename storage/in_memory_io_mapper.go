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
	return self.Memory.GetVertexById(vertexId, vertex)
}

func (self *InMemoryIOMapper) GetEdgeById(edgeId uuid.UUID, edge *graph.Edge) error {
	return self.Memory.GetEdgeById(edgeId, edge)
}

func (self *InMemoryIOMapper) StoreEdge(edge *graph.Edge, success *bool) error {
	system.Logf(">> Request to Edge %s\n", edge)
	*success = false
	if e := self.Memory.StoreElement(edge); e != nil {
		return e
	}
	if e := self.Memory.Index.CreateEdgeIndex(edge); e != nil {
		return e
	}
	*success = true
	return nil
}

func (self *InMemoryIOMapper) RemoveVertex(vertexId uuid.UUID, success *bool) error {
	*success = false
	if e := self.Memory.RemoveElement(vertexId, graph.VERTEX); e != nil {
		return e
	}
	if e := self.Memory.Index.RemoveVertexIndex(vertexId); e != nil {
		return e
	}
	*success = true
	return nil
}

func (self *InMemoryIOMapper) UpdateProperties(element *graph.Element, success *bool) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) RemoveEdge(edge uuid.UUID, succ *bool) error {
	panic("implement me")
}

var _ IOMapper = new(InMemoryIOMapper)
