package storage

import (
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/metadata"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
)

type InMemoryIOMapper struct {
	BackendId string
	Metadata  metadata.Metadata
	Memory    *MemoryStorage
	//Disk     *DiskStorage
}

func NewInMemoryIOMapper(backendId string, metadataAddrs []string) *InMemoryIOMapper {
	return &InMemoryIOMapper{
		BackendId: backendId,
		Metadata:  metadata.NewZkMetadataMapper(metadataAddrs),
		Memory:    NewMemoryStorage(),
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
	system.Logf(">> Request to Store %s\n", edge)
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

func (self *InMemoryIOMapper) GetOutEdges(vertexId uuid.UUID, edges *[]graph.Edge) error {
	var edge graph.Edge
	var _edges []graph.Edge
	i := self.Memory.Index.VertexIndex[vertexId.String()].Front()
	for i != nil{
		u, _ := uuid.Parse(i.Value.(EdgeIndex).Id)

		e := self.GetEdgeById(u, &edge)
		if e != nil{
			continue
		}
		_edges = append(_edges, edge)
		i = i.Next()
	}
	*edges = _edges
	return nil
}

func (self *InMemoryIOMapper) GetInEdges(vertexId uuid.UUID, edges *[]graph.Edge) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) UpdateProperties(element *graph.Element, success *bool) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) RemoveEdge(edge uuid.UUID, succ *bool) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) RegisterToHostPartition(ids []uuid.UUID, succ *bool) error {
	_, e := self.Metadata.AddBackendToPartition(ids[0], ids[1], self.BackendId)
	return e
}

var _ IOMapper = new(InMemoryIOMapper)
