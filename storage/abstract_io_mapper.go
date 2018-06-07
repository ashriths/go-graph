package storage

import "github.com/google/uuid"

// IO Mapper is the Storage Backend Interface
type IOMapper interface {
	StoreVertex(properties interface{}) (error, uuid.UUID)
	StoreEdge(srcVertex uuid.UUID, destVertex uuid.UUID, properties interface{}) (error, uuid.UUID)

	UpdateProperties(elementId uuid.UUID, properties interface{}) error
	UpdatePropertyByName(elementId uuid.UUID, key string, value string) error

	GetElementProperties(elementId uuid.UUID) (error, interface{})
	GetElementPropertyByName(elementId uuid.UUID, key string) (error, interface{})

	RemoveVertex(vertex uuid.UUID) error
	RemoveEdge(edge uuid.UUID) error
}
