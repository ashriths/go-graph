package graph

type EdgeInterface interface {
	GetSrcVertex() VertexInterface
	GetDestVertex() VertexInterface
}

type Edge struct {
	Element
	SrcVertex  VertexInterface
	DestVertex VertexInterf
	ace
	Properties map[string]string
}

func (self *Edge) GetSrcVertex() VertexInterface {
	panic("todo")
}

func (self *Edge) GetDestVertex() VertexInterface {
	panic("todo")
}
