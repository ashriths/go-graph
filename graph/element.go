package graph

import (
	"encoding/json"
	"github.com/ashriths/go-graph/common"
	"github.com/google/uuid"
)

type ElementProperty map[string]string

type ElementInterface interface {
	GetUUID() uuid.UUID
	GetGraphId() uuid.UUID
	GetProperties() (error, ElementProperty)
	SetProperties(props ElementProperty) error
	Remove() error

	Json() string
}

type Element struct {
	UUID       uuid.UUID
	GraphUUID  uuid.UUID
	Properties ElementProperty
}

func (self *Element) GetUUID() uuid.UUID {
	return self.UUID
}

func (self *Element) GetGraphId() uuid.UUID {
	return self.GraphUUID
}

func (self *Element) GetProperties() (error, ElementProperty) {
	panic("todo")
}

func (self *Element) SetProperties(properties ElementProperty) error {
	panic("todo")
}

func (self *Element) Remove() error {
	panic("todo")
}

func (self *Element) String() string {
	str, e := json.Marshal(self)
	common.LogError(e)
	return string(str)
}

func (self *Element) Json() string {
	str, e := json.Marshal(self)
	common.LogError(e)
	return string(str)
}

//func ElementFactory(elemType string, data interface{}) (interface{}, error){
//	switch elemType {
//	case VERTEX_PREFIX:
//		return new(Vertex(data))
//
//	}
//}

var _ ElementInterface = new(Element)
