package storage

import (
	"container/list"
	"fmt"
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/metadata"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

type InMemoryIOMapper struct {
	BackendId string
	Metadata  metadata.Metadata
	Memory    *MemoryStorage
	//Disk     *DiskStorage
}

func NewInMemoryIOMapper(backendId string, metadataAddrs []string) *InMemoryIOMapper {
	return &InMemoryIOMapper{
		BackendId: backendId,
		Metadata:  metadata.NewZkMetadataMapper(metadataAddrs),
		Memory:    NewMemoryStorage(),
	}
}

func (self *InMemoryIOMapper) StoreVertex(vertex *graph.Vertex, success *bool) error {
	system.Logf(">> Request to Store %s\n", vertex)
	*success = false
	if e := self.Memory.StoreElement(vertex); e != nil {
		return e
	}
	if e := self.Memory.Index.CreateVertexIndex(vertex); e != nil {
		return e
	}
	*success = true
	return nil
}

func (self *InMemoryIOMapper) GetVertexById(vertexId uuid.UUID, vertex *graph.Vertex) error {
	return self.Memory.GetVertexById(vertexId, vertex)
}

func (self *InMemoryIOMapper) GetEdgeById(edgeId uuid.UUID, edge *graph.Edge) error {
	return self.Memory.GetEdgeById(edgeId, edge)
}

func (self *InMemoryIOMapper) StoreEdge(edge *graph.Edge, success *bool) error {
	system.Logf(">> Request to Store %s\n", edge)
	*success = false
	if e := self.Memory.StoreElement(edge); e != nil {
		return e
	}
	if e := self.Memory.Index.CreateEdgeIndex(edge); e != nil {
		return e
	}
	*success = true
	return nil
}

func (self *InMemoryIOMapper) RemoveVertex(vertexId uuid.UUID, success *bool) error {
	*success = false
	if e := self.Memory.RemoveElement(vertexId, graph.VERTEX); e != nil {
		return e
	}
	if e := self.Memory.Index.RemoveVertexIndex(vertexId); e != nil {
		return e
	}
	*success = true
	return nil
}

func (self *InMemoryIOMapper) GetOutEdges(vertexId uuid.UUID, edges *[]graph.Edge) error {
	var edge graph.Edge
	var _edges []graph.Edge
	i := self.Memory.Index.VertexOutEdgeIndex[vertexId.String()].Front()
	for i != nil {
		u, _ := uuid.Parse(i.Value.(EdgeIndex).Id)

		e := self.GetEdgeById(u, &edge)
		if e != nil {
			continue
		}
		_edges = append(_edges, edge)
		i = i.Next()
	}
	*edges = _edges
	return nil
}

func (self *InMemoryIOMapper) GetInEdges(vertexId uuid.UUID, edges *[]graph.Edge) error {
	var _edges []graph.Edge
	var edge graph.Edge
	var e error
	var edgeId uuid.UUID
	edgeMap := self.Memory.Index.EdgeIndex
	for k, v := range edgeMap {
		if v == vertexId.String() && strings.HasSuffix(k, "::"+DEST_PREFIX) {
			edgeId, e = uuid.Parse(strings.Split(k, "::")[0])
			if e != nil {
				return e
			}
			e = self.GetEdgeById(edgeId, &edge)
			if e != nil {
				return e
			}
			_edges = append(_edges, edge)
		}
	}
	*edges = _edges
	return nil
}

func (self *InMemoryIOMapper) UpdateProperties(element *graph.Element, success *bool) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) RemoveEdge(edge uuid.UUID, succ *bool) error {
	panic("implement me")
}

func (self *InMemoryIOMapper) UpdateReplica(data interface{}, succ *bool) error {
	self.Memory.kvStore.dataLock.Lock()
	defer self.Memory.kvStore.dataLock.Unlock()
	switch v := data.(type) {
	case map[string]string:
		self.Memory.kvStore.data = v
	case map[string]*list.List:
		self.Memory.kvStore.listData = v
	default:
		*succ = false
		return fmt.Errorf("Invalid data type")
	}
	*succ = true
	return nil
}

func (self *InMemoryIOMapper) RegisterToHostPartition(ids []uuid.UUID, succ *bool) error {
	_, _, e := self.Metadata.AddBackendToPartition(ids[0], ids[1], self.BackendId)
	if e != nil {
		return e
	}
	//go self.startWatchingPartition(ids[0], ids[1], watch)
	return e
}

func (self *InMemoryIOMapper) startWatchingPartition(graphId uuid.UUID, partitionId uuid.UUID, watch interface{}) {
	var _watch <-chan zk.Event
	_watch = watch.(<-chan zk.Event)
	for {
		evt := <-_watch
		system.Logln("Watch fired for ", partitionId.String(), evt.Path, evt.Type)
		if evt.Type == zk.EventNodeDeleted{
			newBack, e := self.Metadata.FindNewBackendForPartition(graphId, partitionId)
			if e != nil {
				system.Logln("Cannot get new partition to replicate", e)
			}
			e = self.replicateToBackend(graphId, partitionId, newBack)
			if e != nil {
				system.Logln("Error while replicating. ", e)
			}
		}
	}
}

func (self *InMemoryIOMapper) replicateToBackend(graphId uuid.UUID, partitionID uuid.UUID, backendId string) error {
	self.Memory.kvStore.dataLock.Lock()
	defer self.Memory.kvStore.dataLock.Unlock()
	var succ bool
	var err error
	info, err := self.Metadata.GetBackendInformation(backendId)
	if err != nil {
		system.Logln("Failed to fetch backend Info")
		return err
	}
	backendAddr := info["address"].(string)
	client := NewStorageClient(backendAddr)
	err = client.UpdateReplica(self.Memory.kvStore.data, &succ)
	if err != nil {
		return err
	}
	if !succ {
		return fmt.Errorf("Unable to update replica")
	}
	err = client.UpdateReplica(self.Memory.kvStore.listData, &succ)
	if err != nil {
		return err
	}
	if !succ {
		return fmt.Errorf("Unable to update replica")
	}
	err = client.RegisterToHostPartition([]uuid.UUID{graphId, partitionID}, &succ)
	if err != nil {
		return err
	}
	if !succ {
		return fmt.Errorf("Unable to make the other replica register to host")
	}
	return nil
}

var _ IOMapper = new(InMemoryIOMapper)
