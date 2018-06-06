package graph

type VertexInterface interface {
	AddEdge(label string, destVertex VertexInterface, properties map[string]string) error
	GetInEdges(edgeLabels []string) []EdgeInterface
	GetOutEdges(edgeLabels []string) []EdgeInterface
	GetParentVertices(edgeLabels []string) []VertexInterface
	GetChildVertices(edgeLabels []string) []VertexInterface
}

type Vertex struct {
	Element
	OutEdges   []EdgeInterface
	InEdges    []EdgeInterface
	Properties map[string]string
}

func (self *Vertex) AddEdge(label string, destVertex VertexInterface, properties map[string]string) error {
	panic("todo")
}

func (self *Vertex) GetInEdges(edgeLabels []string) []EdgeInterface {
	panic("todo")
}

func (self *Vertex) GetOutEdges(edgeLabels []string) []EdgeInterface {
	panic("todo")
}

func (self *Vertex) ParentVertices(edgeLabels []string) []VertexInterface {
	panic("todo")
}

func (self *Vertex) ChildVertices(edgeLabels []string) []VertexInterface {
	panic("todo")
}
