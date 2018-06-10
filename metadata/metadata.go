package metadata

import (
	"github.com/google/uuid"
	"github.com/samuel/go-zookeeper/zk"
)

// Metadata : Interface exposing zookeeper functionality
type Metadata interface {

	//Intitialize the zookeeper tree
	Initialize() error
	//creates a Znode for a graph vertex
	CreateVertex(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error
	//creates a Znode for a graph edge
	CreateEdge(graphID uuid.UUID, edgeID uuid.UUID, partitionID uuid.UUID) error
	//creates a Znode for a graph partition
	CreatePartition(graphID uuid.UUID, partitionID uuid.UUID) error
	//creates a Znode for a backend
	CreateBackend(backendAddr string) (string, error)
	//creates a Znode for a graph
	CreateGraph(graphID uuid.UUID) error

	//returns the backends that house a particular vertex
	GetVertexLocation(graphID uuid.UUID, vertexID uuid.UUID) ([]string, error)
	//sets the partition to which a vertex belongs
	SetVertexLocation(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error
	//returns the backend that houses the source vertex of an edge
	GetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID) ([]string, error)
	//sets the source vertex of an edge
	SetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID, vertexID uuid.UUID) error
	//returns list of backend IDs
	GetAllBackends() ([]string, error)
	//returns backend information
	GetBackendInformation(backendID string) (map[string]interface{}, error)
	//returns all partitions of a graph
	GetAllPartitions(graphID uuid.UUID) ([]string, error)
	//sets the data at a partition Znode
	SetPartitionInformation(graphID uuid.UUID, partitionID uuid.UUID, data interface{}) error
	//returns the data at a partition Znode
	GetPartitionInformation(graphID uuid.UUID, partitionID uuid.UUID) (map[string]interface{}, error)

	//adds a backend to a partition
	AddBackendToPartition(graphID uuid.UUID, partitionID uuid.UUID, backendID string) ([]string, <-chan zk.Event, error)
}
