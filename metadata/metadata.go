package metadata

import (
	"github.com/google/uuid"
)

// Metadata : Interface exposing zookeeper functionality
type Metadata interface {

	//Intitialize the zookeeper tree
	Initialize() error
	//creates a Znode for a graph vertex
	CreateVertex(graphID uuid.UUID, partitionID uuid.UUID, vertexID uuid.UUID) error
	//creates a Znode for a graph edge
	CreateEdge(graphID uuid.UUID, partitionID uuid.UUID, edgeID uuid.UUID) error
	//creates a Znode for a graph partition
	CreatePartition(graphID uuid.UUID, partitionID uuid.UUID) error
	//creates a Znode for a backend
	CreateBackend(backendAddr string) (string, error)

	//creates a Znode for a graph
	CreateGraph(graphID uuid.UUID, data interface{}) error
	// Gets all graphs in the system
	GetGraphs(graphsIds *[]uuid.UUID) error

	//Delete the Znode for a graph vertex
	DeleteVertex(graphID uuid.UUID, vertexID uuid.UUID) error
	//Delete the Znode for a graph edge
	DeleteEdge(graphID uuid.UUID, edgeID uuid.UUID) error

	//returns the backends that house a particular vertex
	GetVertexLocation(graphID uuid.UUID, vertexID uuid.UUID) (*uuid.UUID, []string, error)
	//sets the partition to which a vertex belongs
	SetVertexLocation(graphID uuid.UUID, partitionID uuid.UUID, vertexID uuid.UUID) error
	//returns the backend that houses the source vertex of an edge
	GetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID) (*uuid.UUID, []string, error)
	//sets the source vertex of an edge
	SetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID, vertexID uuid.UUID) error
	//returns list of backend IDs
	GetAllBackends() ([]string, error)
	//return all backends that house this partition
	GetBackendsForPartition(graphID uuid.UUID, partitionID uuid.UUID) ([]string, error)
	//returns backend information
	GetBackendInformation(backendID string) (map[string]interface{}, error)
	//returns all partitions of a graph
	GetAllPartitions(graphID uuid.UUID) ([]string, error)
	//sets the data at a partition Znode
	SetPartitionInformation(graphID uuid.UUID, partitionID uuid.UUID, data interface{}) error
	//returns the data at a partition Znode
	GetPartitionInformation(graphID uuid.UUID, partitionID uuid.UUID) (map[string]interface{}, error)
	//  updates the data at a vertex Znode
	UpdateVertexInformation(graphID uuid.UUID, vertexID uuid.UUID, key interface{}, value interface{}) error
	// returns data at a vertex Znode
	DeleteVertexInformation(graphID uuid.UUID, vertexID uuid.UUID, key interface{}) error
	// returns the data at a vertex Znode
	GetVertexInformation(graphID uuid.UUID, vertexID uuid.UUID) (map[string]interface{}, error)
	//  updates the data at an Edge Znode
	UpdateEdgeInformation(graphID uuid.UUID, edgeID uuid.UUID, key interface{}, value interface{}) error
	// returns the data at a vertex Znode
	GetEdgeInformation(graphID uuid.UUID, edgeID uuid.UUID) (map[string]interface{}, error)

	//adds a backend to a partition
	AddBackendToPartition(graphID uuid.UUID, partitionID uuid.UUID, backendID string) ([]string, interface{}, error)
	//increment count of number of edges/vertices in a partition
	IncrementElementCount(graphID uuid.UUID, partitionID uuid.UUID) error
	//decrement count of number of edges/vertices in a partition
	DecrementElementCount(graphID uuid.UUID, partitionID uuid.UUID) error

	FindNewBackendForPartition(graphID uuid.UUID, partitionID uuid.UUID) (string, error)
}
