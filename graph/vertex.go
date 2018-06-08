package graph

type Vertex interface {
	GetInEdges(edgeLabels []string) (error, []Edge)
	GetOutEdges(edgeLabels []string) (error, []Edge)
	GetParentVertices(edgeLabels []string) (error, []Vertex)
	GetChildVertices(edgeLabels []string) (error, []Vertex)
}

type GoGraphVertex struct {
	GoGraphElement
	OutEdges   []Edge
	InEdges    []Edge
}


func (self *GoGraphVertex) GetInEdges(edgeLabels []string) (error, []Edge) {
	panic("todo")
}

func (self *GoGraphVertex) GetOutEdges(edgeLabels []string) (error,[]Edge) {
	panic("todo")
}

func (self *GoGraphVertex) GetParentVertices(edgeLabels []string) (error, []Vertex) {
	panic("todo")
}

func (self *GoGraphVertex) GetChildVertices(edgeLabels []string) (error, []Vertex) {
	panic("todo")
}

func V(data string) *GoGraphVertex{
	return &GoGraphVertex{
		GoGraphElement: GoGraphElement{
			Label:data,
			UUID:nil,
			Properties:nil,
		},
		OutEdges:nil,
		InEdges:nil,
	}
}

var _ Vertex = new(GoGraphVertex)