package system

import "github.com/google/uuid"

func NewUUID() (error, uuid.UUID) {
	uuid, e := uuid.NewUUID()
	return e, uuid
}
