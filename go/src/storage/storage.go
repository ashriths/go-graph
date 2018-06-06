package storage

import (
	"go-graph/go/src/graph"
	"github.com/google/uuid"
)

type Storage interface {
	GetVertex(uuid uuid.UUID) (error, *graph.Vertex)
	GetEdge(uuid uuid.UUID) (error, *graph.Edge)
	GetEdges(node *graph.Vertex) (error, []*graph.Edge)
	GetEdgeByName(node *graph.Vertex, relationName string) (error, *graph.Edge)
}
