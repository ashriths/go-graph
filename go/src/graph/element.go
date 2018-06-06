package graph

import "github.com/google/uuid"

type Element interface {
	GetUUID() uuid.UUID
	GetLabel() (error, string)
	Graph() (error, Graph)
	GetProperties() (error, interface{})
	SetProperties(props interface{}) error
	Remove() error
}

type GoGraphElement struct {
	UUID  uuid.UUID
	Label string
	Properties interface{}
}

func (self *GoGraphElement) GetUUID() uuid.UUID {
	return self.UUID
}

func (self *GoGraphElement) GetLabel() (error, string) {
	return nil, self.Label
}

func (self *GoGraphElement) Graph() (error, Graph) {
	panic("todo")
}

func (self *GoGraphElement) GetProperties() (error, interface{})  {
	panic("todo")
}

func (self *GoGraphElement) SetProperties(properties interface{}) error  {
	panic("todo")
}

func (self *GoGraphElement) Remove() error {
	panic("todo")
}

var _ Element = new(GoGraphElement)
