package storage

import (
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
	"net/rpc"
	"reflect"
)

type StorageClient struct {
	Addr string
	Conn *rpc.Client
}

func NewStorageClient(addr string) *StorageClient {
	return &StorageClient{Addr: addr, Conn: nil}
}

func (self *StorageClient) Connect(force bool) error {

	// Create a persistent connection to the server
	if force || self.Conn == nil {
		if self.Conn != nil {
			system.Logln("Closing stale connection..")
			self.Conn.Close()
		}
		system.Logln(self.Addr, ">> Creating connection to ", self.Addr)
		conn, e := rpc.DialHTTP("tcp", self.Addr)
		if e != nil {
			return e
		}
		system.Logln(self.Addr, "<< Connection Successful")
		self.Conn = conn
	}
	return nil
}

func (self *StorageClient) Call(method string, args interface{}, reply interface{}) error {
	self.Connect(false)
	system.Logf("%v >> %v Args: %v(%T)\n", self.Addr, method, args, args)
	if e := self.Conn.Call(method, args, reply); e != nil {
		system.Logln("Error", e)
		if e := self.Connect(true); e != nil {
			return e
		}
	}
	system.Logf("%v << %v Args: %v %v\n", self.Addr, method, args, reflect.ValueOf(reply).Elem())
	return nil
}

func (self *StorageClient) StoreVertex(vertex *graph.Vertex, succ *bool) error {
	return self.Call("Storage.StoreVertex", vertex, succ)
}

func (self *StorageClient) StoreEdge(edge *graph.Edge, succ *bool) error {
	return self.Call("Storage.StoreEdge", edge, succ)
}

func (self *StorageClient) GetVertexById(vertexId uuid.UUID, vertex *graph.Vertex) error {
	return self.Call("Storage.GetVertexById", vertexId, vertex)
}

func (self *StorageClient) GetEdgeById(edgeId uuid.UUID, edge *graph.Edge) error {
	return self.Call("Storage.GetEdgeById", edgeId, edge)
}

func (self *StorageClient) UpdateProperties(element *graph.Element, success *bool) error {
	return self.Call("Storage.UpdateProperties", element, success)
}

func (self *StorageClient) RemoveVertex(vertex uuid.UUID, succ *bool) error {
	return self.Call("Storage.RemoveVertex", vertex, succ)
}

func (self *StorageClient) RemoveEdge(edge uuid.UUID, succ *bool) error {
	return self.Call("Storage.RemoveEdge", edge, succ)
}

func (self *StorageClient) RegisterToHostPartition(ids []uuid.UUID, succ *bool) error {
	return self.Call("Storage.RegisterToHostPartition", ids, succ)
}

func (self *StorageClient) GetOutEdges(vertexId uuid.UUID, edges *[]graph.Edge) error {
	return self.Call("Storage.GetOutEdges", vertexId, edges)
}

func (self *StorageClient) GetInEdges(vertexId uuid.UUID, edges *[]graph.Edge) error {
	return self.Call("Storage.GetInEdges", vertexId, edges)
}

func (self *StorageClient) UpdateReplica(data interface{}, succ *bool) error {
	return self.Call("Storage.UpdateReplica", data, succ)
}

var _ IOMapper = new(StorageClient)
