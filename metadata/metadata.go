package metadata

import (
	"github.com/google/uuid"
	"github.com/ashriths/go-graph/storage"
)

// Metadata : Interface exposing zookeeper functionality
type Metadata interface {
	//creates a Znode for a graph vertex
	createVertexZnode(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error
	//creates a Znode for a graph edge
	createEdgeZnode(graphID uuid.UUID, edgeID uuid.UUID, partitionID uuid.UUID) error
<<<<<<< Updated upstream
	//Method to create a Znode
	getVertexLocation(graphID uuid.UUID, vertexID uuid.UUID) (error, []string)
	//Method to create a Znode
=======
	//returns the backends that houses a particular vertex
	getVertexLocation(graphID uuid.UUID, vertexID uuid.UUID) ([]string, error)
	//sets the partition to which a vertex belongs
>>>>>>> Stashed changes
	setVertexLocation(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error
	//returns the backend that houses the source vertex of an edge
	getEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID) ([]string, error)
	//sets the
	setEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID, partitionID uuid.UUID) error
	//Method for a backend to add itself to the zookeeper metadata
	addBackend(backendID uuid.UUID, backendAddr string) error
}
