package storage

import "github.com/google/uuid"

// IO Mapper is the Storage Backend Interface
type IOMapper interface {
	Get(nodeId uuid.UUID)
}
