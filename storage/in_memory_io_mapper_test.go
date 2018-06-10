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
