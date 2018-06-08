package graph

import "github.com/google/uuid"

type Graph interface {
	GetUUID() uuid.UUID
	GetLabel() (error, string)
	GetProperties() (error, interface{})
	GetVertices() (error, []Vertex)
	GetEdges() (error, []Edge)
	AddVertex(properties interface{}) (error, Vertex)
	AddEdge(src Vertex, dest Vertex, properties interface{}) (error, Edge)
}

type GoGraph struct {
	UUID uuid.UUID
	Label string
	Properties interface{}
	Vertices []Vertex
	Edges []Edge
}

func NewGraph(interface{}) (error, uuid.UUID){
	panic("todo")
}

func GetGraph(uuid uuid.UUID) Graph{
	panic("todo")
}

func (self *GoGraph) GetUUID() uuid.UUID {
	return self.UUID
}

func (self *GoGraph) GetLabel() (error, string) {
	return nil, self.Label
}

func (self *GoGraph) GetProperties() (error, interface{}){
	panic("todo")
}

func (self *GoGraph) GetVertices() (error, []Vertex) {
	panic("todo")
}

func (self *GoGraph) GetEdges() (error, []Edge) {
	panic("todo")
}

func (self *GoGraph) AddVertex(properties interface{}) (error, Vertex){
	panic("todo")
}

func (self *GoGraph) AddEdge(src Vertex, dest Vertex, properties interface{}) (error, Edge){
	panic("todo")
}

var _ Graph = new(GoGraph)
