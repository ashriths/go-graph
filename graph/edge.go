package graph

type EdgeInterface interface {
	GetSrcVertex() (error, Vertex)
	GetDestVertex() (error, Vertex)
}

type Edge struct {
	Element
	SrcVertex  Vertex
	DestVertex Vertex
}

func (self *Edge) GetSrcVertex() (error, Vertex) {
	panic("todo")
}

func (self *Edge) GetDestVertex() (error, Vertex) {
	panic("todo")
}

var _ EdgeInterface = new(Edge)
