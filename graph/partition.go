package graph

import "github.com/google/uuid"

type PartitionInterface interface {
	Move(partitionInterface PartitionInterface) error
}

type Partition struct {
	Id       uuid.UUID
	vertices []Vertex
	edges    []Edge
}

func (partition *Partition) Move(partitionInterface PartitionInterface) error {
	panic("implement me")
}
