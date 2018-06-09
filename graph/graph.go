package graph

import (
	"encoding/json"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
)

type GraphInterface interface {
	GetUUID() uuid.UUID
	GetProperties() (error, ElementProperty)
	GetVertices() (error, []Vertex)
	GetEdges() (error, []Edge)
	AddVertex(properties ElementProperty) (error, Vertex)
	AddEdge(src Vertex, dest Vertex, properties ElementProperty) (error, Edge)
}

const (
	GRAPH  = "GRAPH"
	VERTEX = "VERTEX"
	EDGE   = "EDGE"
)

type Graph struct {
	UUID       uuid.UUID
	Properties interface{}
	Vertices   []Vertex
	Edges      []Edge
}

func NewGraph(interface{}) (error, uuid.UUID) {
	panic("todo")
}

func GetGraph(uuid uuid.UUID) Graph {
	panic("todo")
}

func GetElementType(elem interface{}) string {
	switch elem.(type) {
	case Graph:
		return GRAPH
	case *Graph:
		return GRAPH
	case Vertex:
		return VERTEX
	case *Vertex:
		return VERTEX
	case Edge:
		return EDGE
	case *Edge:
		return EDGE
	default:
		panic("Invalid type")
	}
}

func GetElement(elemType, data string, elem interface{}) error {
	var e error
	switch elemType {
	case GRAPH:
		var graph Graph
		e = json.Unmarshal([]byte(data), &graph)
		elem = graph
	case VERTEX:
		var vertex Vertex
		system.Logln("Before Unmarshall", data)
		e = json.Unmarshal([]byte(data), &vertex)
		system.Logln("After Unmarshall", vertex)
		elem = vertex
		system.Logln("After TypeCast", elem)
	case EDGE:
		var edge Edge
		e = json.Unmarshal([]byte(data), &edge)
		elem = edge
	default:
		panic("Invalid type")
	}
	return e
}

func (self *Graph) GetUUID() uuid.UUID {
	return self.UUID
}

func (self *Graph) GetProperties() (error, ElementProperty) {
	panic("todo")
}

func (self *Graph) GetVertices() (error, []Vertex) {
	panic("todo")
}

func (self *Graph) GetEdges() (error, []Edge) {
	panic("todo")
}

func (self *Graph) AddVertex(properties ElementProperty) (error, Vertex) {
	panic("todo")
}

func (self *Graph) AddEdge(src Vertex, dest Vertex, properties ElementProperty) (error, Edge) {
	panic("todo")
}

var _ GraphInterface = new(Graph)
