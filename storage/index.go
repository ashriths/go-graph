package storage

import (
	"container/list"
	"github.com/ashriths/go-graph/graph"
	"sync"
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
	panic("todo")
}
