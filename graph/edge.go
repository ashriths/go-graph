package graph

type Edge interface {
	GetSrcVertex() (error, Vertex)
	GetDestVertex() (error, Vertex)
}

type GoGraphEdge struct {
	Element
	SrcVertex  Vertex
	DestVertex Vertex
}

func (self *GoGraphEdge) GetSrcVertex() (error, Vertex) {
	panic("todo")
}

func (self *GoGraphEdge) GetDestVertex() (error, Vertex) {
	panic("todo")
}

var _ Edge = new(GoGraphEdge)
