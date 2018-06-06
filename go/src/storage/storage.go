package storage

import "github.com/google/uuid"

type Storage interface {
	AddDocument(properties interface{}) (error, uuid.UUID)
	DeleteDocument(uuid uuid.UUID) error
	UpdateDocument(uuid uuid.UUID, properties interface{}) error
}

