package graph

import "github.com/google/uuid"

type GraphInterface interface {
	GetVertex() VertexInterface
	GetEdge() EdgeInterface
}

type Graph struct {
	uuid uuid.UUID
	v []VertexInterface
	e []EdgeInterface
}

func (self *Graph) GetVertex() VertexInterface {
	panic("todo")
}

func (self *Graph) GetEdge() VertexInterface {
	panic("todo")
}
