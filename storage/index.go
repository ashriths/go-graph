package storage

import (
	"container/list"
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
	"sync"
)

const (
	SRC_PREFIX  = "SRC"
	DEST_PREFIX = "DST"
)

type Index struct {
	VertexIndex      map[string]*list.List
	EdgeIndex        map[string]string
	GraphVertexIndex map[string]*list.List
	GraphEdgeIndex   map[string]*list.List
	indexLock        sync.Mutex
}

type EdgeIndex struct {
	Name string
	Id   string
}

func NewIndex() *Index {
	return &Index{
		VertexIndex:      make(map[string]*list.List),
		EdgeIndex:        make(map[string]string),
		GraphVertexIndex: make(map[string]*list.List),
		GraphEdgeIndex:   make(map[string]*list.List),
	}
}

func (self *Index) CreateVertexIndex(vertex *graph.Vertex) error {
	self.indexLock.Lock()
	defer self.indexLock.Unlock()

	self.VertexIndex[vertex.GetUUID().String()] = list.New()
	system.Logln(vertex.Json())

	//_, ok := self.GraphVertexIndex[vertex.GetGraphId().String()]
	//if !ok {
	//	self.GraphVertexIndex[vertex.GetGraphId().String()] = list.New()
	//}
	//self.GraphVertexIndex[vertex.GetGraphId().String()].PushBack(vertex.GetUUID().String())

	return nil
}

func (self *Index) CreateEdgeIndex(edge *graph.Edge) error {
	self.indexLock.Lock()
	defer self.indexLock.Unlock()

	key := edge.GetUUID().String() + "::" + SRC_PREFIX
	self.EdgeIndex[key] = edge.SrcVertex.String()

	self.VertexIndex[edge.SrcVertex.String()].PushBack(EdgeIndex{
		Name: edge.Name,
		Id:   edge.GetUUID().String(),
	})

	key = edge.GetUUID().String() + "::" + DEST_PREFIX
	self.EdgeIndex[key] = edge.DestVertex.String()
	return nil
}

func (self *Index) RemoveVertexIndex(vertexId uuid.UUID) error {
	self.indexLock.Lock()
	defer self.indexLock.Unlock()

	outEdges := self.VertexIndex[vertexId.String()]
	delete(self.VertexIndex, vertexId.String())

	//Remove from EdgeIndex
	i := outEdges.Front()
	for i != nil {
		self.protectedRemoveEdgeIndex(i.Value.(EdgeIndex).Id)
		i = i.Next()
	}
	return nil
}

// This is not synchronized. Call this only if you hold the lock to index
func (self *Index) protectedRemoveEdgeIndex(edgeId string) {
	delete(self.EdgeIndex, edgeId+"::"+DEST_PREFIX)
	delete(self.EdgeIndex, edgeId+"::"+SRC_PREFIX)
}

func (self *Index) RemoveEdgeIndex(edgeId uuid.UUID) error {
	self.indexLock.Lock()
	defer self.indexLock.Unlock()

	self.protectedRemoveEdgeIndex(edgeId.String())
	return nil
}
