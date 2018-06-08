package metadata

import (
	"github.com/google/uuid"
)

// Metadata : Interface exposing zookeeper functionality
type Metadata interface {
	//creates a Znode for a graph vertex
	createVertexZnode(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error
	//creates a Znode for a graph edge
	createEdgeZnode(graphID uuid.UUID, edgeID uuid.UUID, partitionID uuid.UUID) error
	//Method to create a Znode
	getVertexLocation(graphID uuid.UUID, vertexID uuid.UUID) (error, []string)
	//Method to create a Znode
	setVertexLocation(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error
	//Method to create a Znode
	getEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID) (error, []string)
	//Method to create a Znode
	setEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID, partitionID uuid.UUID) error
}
