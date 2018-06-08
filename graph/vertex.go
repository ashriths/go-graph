package graph

type VertexInterface interface {
	GetInEdges(edgeLabels []string) (error, []Edge)
	GetOutEdges(edgeLabels []string) (error, []Edge)
	GetParentVertices(edgeLabels []string) (error, []Vertex)
	GetChildVertices(edgeLabels []string) (error, []Vertex)
}

type Vertex struct {
	Element
	OutEdges   []Edge
	InEdges    []Edge
}


func (self *Vertex) GetInEdges(edgeLabels []string) (error, []Edge) {
	panic("todo")
}

func (self *Vertex) GetOutEdges(edgeLabels []string) (error,[]Edge) {
	panic("todo")
}

func (self *Vertex) GetParentVertices(edgeLabels []string) (error, []Vertex) {
	panic("todo")
}

func (self *Vertex) GetChildVertices(edgeLabels []string) (error, []Vertex) {
	panic("todo")
}

func V(data string) *Vertex{
	return &Vertex{
		Element: Element{
			Label:data,
			UUID:nil,
			Properties:nil,
		},
		OutEdges:nil,
		InEdges:nil,
	}
}

var _ VertexInterface = new(Vertex)