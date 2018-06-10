package storage

import (
	"container/list"
	"github.com/ashriths/go-graph/graph"
	"sync"
	"github.com/google/uuid"
)

const (
	SRC_PREFIX = "SRC"
	DEST_PREFIX = "DST"
)

type Index struct {
	VertexIndex map[string]*list.List
	EdgeIndex   map[string]string
	indexLock   sync.Mutex
}

func NewIndex() *Index {
	return &Index{
		VertexIndex: make(map[string]*list.List),
		EdgeIndex:   make(map[string]string),
	}
}

func (self *Index) CreateVertexIndex(vertex *graph.Vertex) error {
	self.indexLock.Lock()
	defer self.indexLock.Unlock()

	self.VertexIndex[vertex.GetUUID().String()] = list.New()
	return nil
}

func (self *Index) CreateEdgeIndex(edge *graph.Edge) error {
	self.indexLock.Lock()
	defer self.indexLock.Unlock()

	key := edge.GetUUID().String() + "::" + SRC_PREFIX
	self.EdgeIndex[key] = edge.SrcVertex.String()

	key = edge.GetUUID().String() + "::" + DEST_PREFIX
	self.EdgeIndex[key] = edge.DestVertex.String()
	return nil
}

func (self *Index) RemoveVertexIndex(vertexId uuid.UUID) error {
	self.indexLock.Lock()
	defer self.indexLock.Unlock()

	delete(self.VertexIndex, vertexId.String())
	return nil
}
