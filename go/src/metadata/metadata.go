package metadata

import (
	"go-graph/go/src/graph"
	"go-graph/go/src/storage"
)

type Metadata interface {
	GetNodeLocation(node *graph.Node) storage.Storage
}
