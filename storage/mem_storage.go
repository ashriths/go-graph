package storage

import (
	"encoding/json"
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
)

type MemoryStorage struct {
	kvStore *MemKVStore
	Index   *Index
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		kvStore: NewMemKVStore(),
		Index:   NewIndex(),
	}
}

func (self *MemoryStorage) StoreElement(element graph.ElementInterface) error {
	elemType := graph.GetElementType(element)
	key := escapeKey(elemType) + "::" + escapeKey(element.GetUUID().String())
	if e := self.kvStore.Set(key, element.Json()); e != nil {
		return e
	}
	return nil
}

func (self *MemoryStorage) GetVertexById(elementId uuid.UUID, vertex *graph.Vertex) error {
	key := escapeKey(graph.VERTEX) + "::" + escapeKey(elementId.String())
	val := self.kvStore.Get(key)
	if e := json.Unmarshal([]byte(val), vertex); e != nil {
		return e
	}
	system.Logln(vertex)
	return nil
}
