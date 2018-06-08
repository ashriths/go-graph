package graph

import "github.com/google/uuid"

type GraphInterface interface {
	GetUUID() uuid.UUID
	GetLabel() (error, string)
	GetProperties() (error, interface{})
	GetVertices() (error, []Vertex)
	GetEdges() (error, []Edge)
	AddVertex(properties interface{}) (error, Vertex)
	AddEdge(src Vertex, dest Vertex, properties interface{}) (error, Edge)
}

type Graph struct {
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

func (self *Graph) GetUUID() uuid.UUID {
	return self.UUID
}

func (self *Graph) GetLabel() (error, string) {
	return nil, self.Label
}

func (self *Graph) GetProperties() (error, interface{}){
	panic("todo")
}

func (self *Graph) GetVertices() (error, []Vertex) {
	panic("todo")
}

func (self *Graph) GetEdges() (error, []Edge) {
	panic("todo")
}

func (self *Graph) AddVertex(properties interface{}) (error, Vertex){
	panic("todo")
}

func (self *Graph) AddEdge(src Vertex, dest Vertex, properties interface{}) (error, Edge){
	panic("todo")
}

var _ GraphInterface = new(Graph)
