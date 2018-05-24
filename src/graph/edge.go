package graph

type Edge interface {
	GetSrcVertex() Vertex
	GetDestVertex() Vertex
}

type GoGraphEdge struct {
	GoGraphElement
	srcVertex  Vertex
	destVertex Vertex
	properties map[string]string
}

func (self *GoGraphEdge) GetSrcVertex() Vertex {
	panic("todo")
}

func (self *GoGraphEdge) GetDestVertex() Vertex {
	panic("todo")
}
