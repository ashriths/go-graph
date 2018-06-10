package storage

import (
	"encoding/json"
	"fmt"
	"github.com/ashriths/go-graph/graph"
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

func (self *MemoryStorage) RemoveElement(elementId uuid.UUID, elemType string) error {
	key := escapeKey(elemType) + "::" + escapeKey(elementId.String())
	if e := self.kvStore.Set(key, ""); e != nil {
		return e
	}
	return nil
}

func (self *MemoryStorage) GetVertexById(elementId uuid.UUID, vertex *graph.Vertex) error {
	key := escapeKey(graph.VERTEX) + "::" + escapeKey(elementId.String())
	val := self.kvStore.Get(key)
	if val == "" {
		return fmt.Errorf("Vertex with id %s doesn't exist.", elementId.String())
	}
	if e := json.Unmarshal([]byte(val), vertex); e != nil {
		return e
	}
	return nil
}

func (self *MemoryStorage) GetEdgeById(elementId uuid.UUID, edge *graph.Edge) error {
	key := escapeKey(graph.EDGE) + "::" + escapeKey(elementId.String())
	val := self.kvStore.Get(key)
	if e := json.Unmarshal([]byte(val), edge); e != nil {
		return e
	}
	return nil
}
