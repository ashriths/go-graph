package metadata

import (
	"go-graph/go/src/storage"
	"github.com/google/uuid"
)

type Metadata interface {
	GetNodeLocation(nodeId uuid.UUID) storage.Storage
}
