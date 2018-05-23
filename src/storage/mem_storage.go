package storage

import "go-graph/src/graph"

type MemStorage struct {
}

func NewMemStorage() *MemStorage {
	return &MemStorage{}
}

func (*MemStorage) GetPath(src *graph.Node, dest *graph.Node) (error, []*graph.Node) {
	panic("todo")
}

func (*MemStorage) GetAllRelations(node *graph.Node) (error, []*graph.Node) {
	panic("todo")
}

func (*MemStorage) GetRelation(node *graph.Node, relationName string) (error, []*graph.Node) {
	panic("todo")
}

var _ Storage = new(MemStorage)
