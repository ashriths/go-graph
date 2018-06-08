package storage

import (
	"github.com/google/uuid"
	"github.com/ashriths/go-graph/graph"
)

// IO Mapper is the Storage Backend Interface
type IOMapper interface {
	StoreVertex(vertex graph.Vertex, uuid *uuid.UUID) error
	StoreEdge(edge *graph.Edge, uuid *uuid.UUID) error

	UpdateProperties(element *graph.Element, success *bool) error

	GetElementProperties(elementId uuid.UUID, properties *interface{}) error

	RemoveVertex(vertex uuid.UUID, succ *bool) error
	RemoveEdge(edge uuid.UUID, succ *bool) error
}
