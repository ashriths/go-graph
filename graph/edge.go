package graph

import (
	"github.com/google/uuid"
	"encoding/json"
	"github.com/ashriths/go-graph/common"
)

type EdgeInterface interface {
	ElementInterface
	GetSrcVertex() (error, Vertex)
	GetDestVertex() (error, Vertex)
}

type Edge struct {
	Element
	Name       string
	SrcVertex  uuid.UUID
	DestVertex uuid.UUID
}

func (self *Edge) GetSrcVertex() (error, Vertex) {
	panic("todo")
}

func (self *Edge) GetDestVertex() (error, Vertex) {
	panic("todo")
}

func (self *Edge) Json() string {
	str, e := json.Marshal(self)
	common.LogError(e)
	return string(str)
}

var _ EdgeInterface = new(Edge)
