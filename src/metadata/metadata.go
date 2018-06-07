package metadata

import (
	"go-graph/src/storage"
	"github.com/google/uuid"
)

type Metadata interface {
	GetNodeLocation(nodeId uuid.UUID) storage.Storage
}
