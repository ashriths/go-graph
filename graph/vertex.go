package graph

import (
	"fmt"
	"github.com/google/uuid"
	"encoding/json"
	"github.com/ashriths/go-graph/common"
)

type VertexInterface interface {
	ElementInterface
	GetInEdges(edgeLabels []string) (error, []Edge)
	GetOutEdges(edgeLabels []string) (error, []Edge)
	GetParentVertices(edgeLabels []string) (error, []Vertex)
	GetChildVertices(edgeLabels []string) (error, []Vertex)
}

type Vertex struct {
	Element
	OutEdges []Edge
	InEdges  []Edge
}

func (self *Vertex) GetInEdges(edgeLabels []string) (error, []Edge) {
	panic("todo")
}

func (self *Vertex) GetOutEdges(edgeLabels []string) (error, []Edge) {
	panic("todo")
}

func (self *Vertex) GetParentVertices(edgeLabels []string) (error, []Vertex) {
	panic("todo")
}

func (self *Vertex) GetChildVertices(edgeLabels []string) (error, []Vertex) {
	panic("todo")
}

func (self *Vertex) String() string {
	return fmt.Sprintf("<Vertex:%s>%s ", self.GetUUID(), self.Element.String())
}

func (self *Vertex) Json() string {
	str, e := json.Marshal(self)
	common.LogError(e)
	return string(str)
}

func V(uuid uuid.UUID, property ElementProperty) *Vertex {
	return &Vertex{
		Element: Element{
			UUID:       uuid,
			Properties: property,
		},
	}
}

func E(uid uuid.UUID, src uuid.UUID, dest uuid.UUID, property ElementProperty) *Edge {
	return &Edge{
		SrcVertex:  src,
		DestVertex: dest,
		Element: Element{
			UUID:       uid,
			Properties: property,
		},
	}
}

var _ VertexInterface = new(Vertex)
