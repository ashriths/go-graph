package metadata

import (
	"github.com/google/uuid"
)

// Metadata : Interface exposing zookeeper functionality
type Metadata interface {
	//creates a Znode for a graph vertex
	CreateVertex(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error
	//creates a Znode for a graph edge
	CreateEdge(graphID uuid.UUID, edgeID uuid.UUID, partitionID uuid.UUID) error

	//returns the backends that houses a particular vertex
	GetVertexLocation(graphID uuid.UUID, vertexID uuid.UUID) ([]string, error)

	//sets the partition to which a vertex belongs
	SetVertexLocation(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error
	//returns the backend that houses the source vertex of an edge
	GetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID) ([]string, error)
	//sets the
	SetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID, partitionID uuid.UUID) error
	//Method for a backend to add itself to the zookeeper metadata
	AddBackend(backendAddr string) error
}
