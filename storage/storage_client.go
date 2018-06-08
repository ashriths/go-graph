package storage

import (
	"net/rpc"
	"reflect"
	"github.com/ashriths/go-graph/graph"
	"github.com/google/uuid"
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
			logln("Closing stale connection..")
			self.Conn.Close()
		}
		logln(self.Addr, ">> Creating connection to ", self.Addr)
		conn, e := rpc.DialHTTP("tcp", self.Addr)
		if e != nil {
			return e
		}
		logln(self.Addr, "<< Connection Successful")
		self.Conn = conn
	}
	return nil
}

func (self *StorageClient) Call(method string, args interface{}, reply interface{}) error {
	self.Connect(false)
	logf("%v >> %v Args: %v(%T)\n", self.Addr, method, args, args)
	if e := self.Conn.Call(method, args, reply); e != nil {
		logln("Error", e)
		if e := self.Connect(true); e != nil {
			return e
		}
		if e := self.Conn.Call(method, args, reply); e != nil {
			return e
		}
	}
	logf("%v << %v Args: %v %v\n", self.Addr, method, args, reflect.ValueOf(reply).Elem())
	return nil
}

func (self *StorageClient) StoreVertex(vertex *graph.Vertex, uuid *uuid.UUID) error {
	return self.Call("Storage.StoreVertex", vertex, uuid)
}

func (self *StorageClient) StoreEdge(edge *graph.Edge, uuid *uuid.UUID) error {
	return self.Call("Storage.StoreEdge", edge, uuid)
}

func (self *StorageClient) UpdateProperties(element *graph.Element, success *bool) error {
	return self.Call("Storage.UpdateProperties", element, success)
}

func (self *StorageClient) GetElementProperties(elementId uuid.UUID, properties *interface{}) error {
	return self.Call("Storage.GetElementProperties", elementId, properties)
}

func (self *StorageClient) RemoveVertex(vertex uuid.UUID, succ *bool) error {
	return self.Call("Storage.RemoveVertex", vertex, succ)
}

func (self *StorageClient) RemoveEdge(edge uuid.UUID, succ *bool) error {
	return self.Call("Storage.RemoveEdge", edge, succ)
}

var _ IOMapper = new(StorageClient)