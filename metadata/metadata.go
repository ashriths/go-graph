package metadata

import (
	"github.com/google/uuid"
	"github.com/samuel/go-zookeeper/zk"
)

// Metadata : Interface exposing zookeeper functionality
type Metadata interface {
	//creates a Znode for a graph vertex
	CreateVertex(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error
	//creates a Znode for a graph edge
	CreateEdge(graphID uuid.UUID, edgeID uuid.UUID, partitionID uuid.UUID) error
	//creates a Znode for a graph partition
	CreatePartition(graphID uuid.UUID, partitionID uuid.UUID) error
	//creates a Znode for a backend
	CreateBackend(backendAddr string) (string, error)

	//returns the backends that house a particular vertex
	GetVertexLocation(graphID uuid.UUID, vertexID uuid.UUID) ([]string, error)
	//sets the partition to which a vertex belongs
	SetVertexLocation(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error
	//returns the backend that houses the source vertex of an edge
	GetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID) ([]string, error)
	//sets the source vertex of an edge
	SetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID, vertexID uuid.UUID) error

	//adds a backend to a partition
	AddBackendToPartition(graphID uuid.UUID, partitionID uuid.UUID, backendID string) ([]string, <-chan zk.Event, error)
}
