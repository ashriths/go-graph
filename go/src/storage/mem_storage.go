package storage

import (
	"go-graph/go/src/graph"
	"github.com/google/uuid"
)

type MemStorage struct {
}

func NewMemStorage() *MemStorage {
	return &MemStorage{}
}

func (*MemStorage) GetVertex(uuid uuid.UUID) (error, *graph.Vertex) {
	panic("todo")
}

func (*MemStorage) GetEdge(uuid uuid.UUID) (error, *graph.Edge) {
	panic("todo")
}

func (*MemStorage) GetEdges(node *graph.Vertex) (error, []*graph.Edge) {
	panic("todo")
}

func (*MemStorage) GetEdgeByName(node *graph.Vertex, relationName string) (error, *graph.Edge) {
	panic("todo")
}

var _ Storage = new(MemStorage)
