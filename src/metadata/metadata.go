package metadata

import (
	"go-graph/src/graph"
	"go-graph/src/storage"
)

type Metadata interface {
	GetNodeLocation(node *graph.Node) storage.Storage
}
