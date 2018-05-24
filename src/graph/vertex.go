package graph

type Vertex interface {
	AddEdge(label string, destVertex Vertex, properties map[string]string) error
	InEdges(edgeLabels []string) []Edge
	OutEdges(edgeLabels []string) []Edge
	ParentVertices(edgeLabels []string) []Vertex
	ChildVertices(edgeLabels []string) []Vertex
}

type GoGraphVertex struct {
	GoGraphElement
	OutEdges   []Edge
	InEdges    []Edge
	properties map[string]string
}

func (self *GoGraphVertex) AddEdge(label string, destVertex Vertex, properties map[string]string) error {
	panic("todo")
}

func (self *GoGraphVertex) InEdges(edgeLabels []string) []Edge {
	panic("todo")
}

func (self *GoGraphVertex) OutEdges(edgeLabels []string) []Edge {
	panic("todo")
}

func (self *GoGraphVertex) ParentVertices(edgeLabels []string) []Vertex {
	panic("todo")
}

func (self *GoGraphVertex) ChildVertices(edgeLabels []string) []Vertex {
	panic("todo")
}
