package graph

type Element interface {
	Id() int64
	Label() string
	Graph() Graph
	GetKeys() map[string]bool
	GetValue(key string) string
	SetKey(key string, value string) bool
	Remove() bool
}

type GoGraphElement struct {
	id    int64
	label string
}

func (self *GoGraphElement) Id() int64 {
	panic("todo")
}

func (self *GoGraphElement) Label() string {
	panic("todo")
}

func (self *GoGraphElement) Graph() Graph {
	panic("todo")
}

func (self *GoGraphElement) GetKeys() map[string]bool {
	panic("todo")
}

func (self *GoGraphElement) GetValue(key string) string {
	panic("todo")
}

func (self *GoGraphElement) SetKey(key string, value string) bool {
	panic("todo")
}

func (self *GoGraphElement) Remove() bool {
	panic("todo")
}
