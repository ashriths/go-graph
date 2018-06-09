package storage

import "github.com/ashriths/go-graph/graph"

type Storage interface {
	StoreElement(element graph.ElementInterface) error
}
