package graph

type Edge interface {
	GetSrcVertex() Vertex
	GetDestVertex() Vertex
}

type GoGraphEdge struct {
	srcVertex  Vertex
	destVertex Vertex
}
