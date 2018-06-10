package storage_test

import (
	"testing"
	"github.com/ashriths/go-graph/storage"
	"github.com/ashriths/go-graph/graph"
	"github.com/google/uuid"
	"github.com/ashriths/go-graph/common"
	"reflect"
)

func newUUID() uuid.UUID{
	u,e:= uuid.NewUUID()
	common.NoError(e)
	return u
}

func TestInMemoryIOMapper_StoreVertex(t *testing.T) {
	s := storage.NewInMemoryIOMapper()

	u := newUUID()
	data := graph.ElementProperty{"key":"value"}

	var succ bool
	e := s.StoreVertex(graph.V(u, data), &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)
}

func TestInMemoryIOMapper_GetVertexById(t *testing.T) {
	s := storage.NewInMemoryIOMapper()

	u := newUUID()
	data := graph.ElementProperty{"key":"value"}

	var succ bool
	e := s.StoreVertex(graph.V(u, data), &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	var v graph.Vertex
	e = s.GetVertexById(u, &v)
	common.Assert(e == nil, t)
	common.Assert(v.GetUUID() == u, t)
	common.Assert(reflect.DeepEqual(v.Properties, data), t)
}

func TestInMemoryIOMapper_RemoveVertex(t *testing.T) {
	s := storage.NewInMemoryIOMapper()

	u := newUUID()
	data := graph.ElementProperty{"key":"value"}

	var succ bool
	e := s.StoreVertex(graph.V(u, data), &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	e = s.RemoveVertex(u, &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	var v graph.Vertex
	e = s.GetVertexById(u, &v)
	common.Assert(e != nil, t)
}

func TestInMemoryIOMapper_StoreEdge(t *testing.T) {
	var succ bool
	var e error
	var ed *graph.Edge

	s := storage.NewInMemoryIOMapper()

	u1 := newUUID()
	data := graph.ElementProperty{"key1":"value1"}
	e = s.StoreVertex(graph.V(u1, data), &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	u2 := newUUID()
	data = graph.ElementProperty{"key2":"value2"}
	e = s.StoreVertex(graph.V(u2, data), &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	u3 := newUUID()
	data = graph.ElementProperty{"key3":"value3"}
	ed = graph.E(u3, u1, u2, data)
	e = s.StoreEdge(ed, &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	e = s.GetEdgeById(u3, ed)
	common.Assert(e == nil, t)
	common.Assert(ed.GetUUID() == u3, t)
	common.Assert(reflect.DeepEqual(ed.Properties, data), t)
}
