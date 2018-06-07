package metadata

import (
	"go-graph/storage"
	"github.com/google/uuid"
)

type Metadata interface {
	GetNodeLocation(nodeId uuid.UUID) storage.Storage
}
