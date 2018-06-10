package storage

import (
	"github.com/ashriths/go-graph/graph"
	"github.com/google/uuid"
)

// IO Mapper is the Storage Backend Interface
type IOMapper interface {
	StoreVertex(vertex *graph.Vertex, success *bool) error
	StoreEdge(edge *graph.Edge, success *bool) error

	GetVertexById(vertexId uuid.UUID, vertex *graph.Vertex) error
	GetEdgeById(edgeId uuid.UUID, edge *graph.Edge) error

	UpdateProperties(element *graph.Element, success *bool) error

	RemoveVertex(vertex uuid.UUID, succ *bool) error
	RemoveEdge(edge uuid.UUID, succ *bool) error

	RegisterToHostPartition(ids []uuid.UUID, succ *bool) error
}
