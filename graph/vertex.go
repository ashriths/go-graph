package graph

import (
	"encoding/json"
	"fmt"
	"github.com/ashriths/go-graph/common"
	"github.com/google/uuid"
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

func V(guuid, uuid uuid.UUID, property ElementProperty) *Vertex {
	return &Vertex{
		Element: Element{
			UUID:       uuid,
			GraphUUID:  guuid,
			Properties: property,
		},
	}
}

var _ VertexInterface = new(Vertex)
