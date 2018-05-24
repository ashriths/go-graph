package graph

type VertexInterface interface {
	AddEdge(label string, destVertex VertexInterface, properties map[string]string) error
	InEdges(edgeLabels []string) []EdgeInterface
	OutEdges(edgeLabels []string) []EdgeInterface
	ParentVertices(edgeLabels []string) []VertexInterface
	ChildVertices(edgeLabels []string) []VertexInterface
}

type Vertex struct {
	Element
	OutEdges   []EdgeInterface
	InEdges    []EdgeInterface
	properties map[string]string
}

func (self *Vertex) AddEdge(label string, destVertex VertexInterface, properties map[string]string) error {
	panic("todo")
}

func (self *Vertex) InEdges(edgeLabels []string) []EdgeInterface {
	panic("todo")
}

func (self *Vertex) OutEdges(edgeLabels []string) []EdgeInterface {
	panic("todo")
}

func (self *Vertex) ParentVertices(edgeLabels []string) []VertexInterface {
	panic("todo")
}

func (self *Vertex) ChildVertices(edgeLabels []string) []VertexInterface {
	panic("todo")
}
