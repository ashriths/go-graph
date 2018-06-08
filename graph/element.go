package graph

import (
	"github.com/google/uuid"
	"encoding/json"
	"github.com/ashriths/go-graph/common"
)

type ElementInterface interface {
	GetUUID() uuid.UUID
	GetLabel() (error, string)
	Graph() (error, Graph)
	GetProperties() (error, interface{})
	SetProperties(props interface{}) error
	Remove() error
}

type Element struct {
	UUID  *uuid.UUID
	Label string
	Properties interface{}
}

func (self *Element) GetUUID() uuid.UUID {
	return *self.UUID
}

func (self *Element) GetLabel() (error, string) {
	return nil, self.Label
}

func (self *Element) Graph() (error, Graph) {
	panic("todo")
}

func (self *Element) GetProperties() (error, interface{})  {
	panic("todo")
}

func (self *Element) SetProperties(properties interface{}) error  {
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

var _ ElementInterface = new(Element)
