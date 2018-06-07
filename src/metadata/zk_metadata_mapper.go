package metadata

import (
	"go-graph/src/storage"
	"github.com/google/uuid"
)

type ZkMetadataMapper struct {
	Addrs []string
}

func (self *ZkMetadataMapper) GetNodeLocation(nodeId uuid.UUID) storage.Storage {
	panic("todo")
}

var _ Metadata = new(ZkMetadataMapper)
