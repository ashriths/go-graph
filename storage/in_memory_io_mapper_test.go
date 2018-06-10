package storage_test

import (
	"github.com/ashriths/go-graph/common"
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/storage"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func newUUID() uuid.UUID {
	u, e := uuid.NewUUID()
	common.NoError(e)
	return u
}

func TestInMemoryIOMapper_StoreVertex(t *testing.T) {
	s := storage.NewInMemoryIOMapper()

	gu := newUUID()
	u := newUUID()
	data := graph.ElementProperty{"key": "value"}

	var succ bool
	e := s.StoreVertex(graph.V(gu, u, data), &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	// Check Indexes
	v, ok := s.Memory.Index.VertexIndex[u.String()]
	common.Assert(ok == true, t)
	common.Assert(v.Len() == 0, t)
}

func TestInMemoryIOMapper_GetVertexById(t *testing.T) {
	s := storage.NewInMemoryIOMapper()
	gu := newUUID()
	u := newUUID()
	data := graph.ElementProperty{"key": "value"}

	var succ bool
	e := s.StoreVertex(graph.V(gu, u, data), &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	var v graph.Vertex
	e = s.GetVertexById(u, &v)
	common.Assert(e == nil, t)
	common.Assert(v.GetUUID() == u, t)
	common.Assert(reflect.DeepEqual(v.Properties, data), t)
}

func TestInMemoryIOMapper_RemoveVertex(t *testing.T) {
	var succ bool
	var ed *graph.Edge

	s := storage.NewInMemoryIOMapper()
	gu := newUUID()

	u := newUUID()
	data := graph.ElementProperty{"key": "value"}

	e := s.StoreVertex(graph.V(gu, u, data), &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	u2 := newUUID()
	data = graph.ElementProperty{"key2": "value2"}
	e = s.StoreVertex(graph.V(gu, u2, data), &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	u3 := newUUID()
	data = graph.ElementProperty{"key3": "value3"}
	ed = graph.E(gu, u3, u, u2, data)
	e = s.StoreEdge(ed, &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	e = s.GetEdgeById(u3, ed)
	common.Assert(e == nil, t)
	common.Assert(ed.GetUUID() == u3, t)
	common.Assert(reflect.DeepEqual(ed.Properties, data), t)

	e = s.RemoveVertex(u, &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	var v graph.Vertex
	e = s.GetVertexById(u, &v)
	common.Assert(e != nil, t)

	// Check Index
	_, ok := s.Memory.Index.VertexIndex[u.String()]
	common.Assert(ok == false, t)

	// Check for references in Index
	for _, val := range s.Memory.Index.EdgeIndex {
		common.Assert(val != u.String(), t)
	}
}

func TestInMemoryIOMapper_StoreEdge(t *testing.T) {
	var succ bool
	var e error
	var ed *graph.Edge

	s := storage.NewInMemoryIOMapper()
	gu := newUUID()

	u1 := newUUID()
	data := graph.ElementProperty{"key1": "value1"}
	e = s.StoreVertex(graph.V(gu, u1, data), &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	u2 := newUUID()
	data = graph.ElementProperty{"key2": "value2"}
	e = s.StoreVertex(graph.V(gu, u2, data), &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	u3 := newUUID()
	data = graph.ElementProperty{"key3": "value3"}
	ed = graph.E(gu, u3, u1, u2, data)
	e = s.StoreEdge(ed, &succ)
	common.Assert(e == nil, t)
	common.Assert(succ == true, t)

	e = s.GetEdgeById(u3, ed)
	common.Assert(e == nil, t)
	common.Assert(ed.GetUUID() == u3, t)
	common.Assert(reflect.DeepEqual(ed.Properties, data), t)
}
