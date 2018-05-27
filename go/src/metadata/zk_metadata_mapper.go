package metadata

import (
	"go-graph/go/src/graph"
	"go-graph/go/src/storage"
)

type ZkMetadataMapper struct {
	Addrs []string
}

func (self *ZkMetadataMapper) GetNodeLocation(node *graph.Node) storage.Storage {
	panic("todo")
}

var _ Metadata = new(ZkMetadataMapper)
