package storage

import (
	"go-graph/go/src/graph"
	"github.com/google/uuid"
)

type DiskStorage struct {
}

func NewDiskStorage() *DiskStorage {
	return &DiskStorage{}
}

func (*DiskStorage) GetVertex(uuid uuid.UUID) (error, *graph.Vertex) {
	panic("todo")
}

func (*DiskStorage) GetEdge(uuid uuid.UUID) (error, *graph.Edge) {
	panic("todo")
}

func (*DiskStorage) GetEdges(node *graph.Vertex) (error, []*graph.Edge) {
	panic("todo")
}

func (*DiskStorage) GetEdgeByName(node *graph.Vertex, relationName string) (error, *graph.Edge) {
	panic("todo")
}

var _ Storage = new(DiskStorage)
