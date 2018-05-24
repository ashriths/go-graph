package graph

type ElementInterface interface {
	Id() int64
	Label() string
	Graph() Graph
	GetKeys() map[string]bool
	GetValue(key string) string
	SetKey(key string, value string) bool
	Remove() bool
}

type Element struct {
	id    int64
	label string
}

func (self *Element) Id() int64 {
	panic("todo")
}

func (self *Element) Label() string {
	panic("todo")
}

func (self *Element) Graph() Graph {
	panic("todo")
}

func (self *Element) GetKeys() map[string]bool {
	panic("todo")
}

func (self *Element) GetValue(key string) string {
	panic("todo")
}

func (self *Element) SetKey(key string, value string) bool {
	panic("todo")
}

func (self *Element) Remove() bool {
	panic("todo")
}
