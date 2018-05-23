package storage

import "go-graph/src/graph"

type DiskStorage struct {
}

func NewDiskStorage() *DiskStorage {
	return &DiskStorage{}
}

func (*DiskStorage) GetPath(src *graph.Node, dest *graph.Node) (error, []*graph.Node) {
	panic("todo")
}

func (*DiskStorage) GetAllRelations(node *graph.Node) (error, []*graph.Node) {
	panic("todo")
}

func (*DiskStorage) GetRelation(node *graph.Node, relationName string) (error, []*graph.Node) {
	panic("todo")
}

var _ Storage = new(DiskStorage)
