package storage

import (
	"go-graph/go/src/graph"
)

type Storage interface {
	GetPath(src *graph.Node, dest *graph.Node) (error, []*graph.Node)
	GetAllRelations(node *graph.Node) (error, []*graph.Node)
	GetRelation(node *graph.Node, relationName string) (error, []*graph.Node)
}
