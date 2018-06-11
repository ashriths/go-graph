package graph

import (
	"encoding/json"
	"github.com/ashriths/go-graph/common"
	"github.com/google/uuid"
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

func E(guuid, uid, src, dest uuid.UUID, name string, property ElementProperty) *Edge {
	return &Edge{
		SrcVertex:  src,
		DestVertex: dest,
		Name:       name,
		Element: Element{
			GraphUUID:  guuid,
			UUID:       uid,
			Properties: property,
		},
	}
}

var _ EdgeInterface = new(Edge)
