package server

import (
	"encoding/json"
	"fmt"
	"github.com/ashriths/go-graph/common"
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/metadata"
	"github.com/ashriths/go-graph/storage"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
	"log"
	"net/http"
	"reflect"
	"github.com/ashriths/go-graph/locator"
)

type Server struct {
	Config         *ServerConfig
	Metadata       metadata.Metadata
	storageClients map[string]*storage.StorageClient
	Locator		   locator.Locator
}

type ServerConfig struct {
	MetadataServers []string
	Addr            string
	Ready           chan<- bool
}

const (
	LOCATOR = "RandomLocator"
)

func NewZookeeperServer(config *ServerConfig) (error, *Server) {
	zkConnMap := &metadata.ZkMetadataMapper{ZkAddrs: config.MetadataServers}
	locator := &locator.RandomLocator{Metadata:zkConnMap}
	return nil, &Server{Config: config, Metadata: zkConnMap, Locator:locator}
}

func (server *Server) ZkCall(method string, args ...interface{}) []interface{} {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	out := reflect.ValueOf(server.Metadata).MethodByName(method).Call(inputs)

	var output = make([]interface{}, len(out))
	for _, outp := range out[:len(out)] {
		output = append(output, outp.Interface())
	}

	return output
}

func (server *Server) getOrCreateStorageClient(backendAddr string) (*storage.StorageClient, error) {
	// Checks if a valid storage client exists for address. If there isn't it creates and returns one
	if stClient, ok := server.storageClients[backendAddr]; ok {
		return stClient, nil
	}
	server.storageClients[backendAddr] = storage.NewStorageClient(backendAddr)
	stClient, _ := server.storageClients[backendAddr]
	return stClient, nil
}

func (server *Server) addvertex(w http.ResponseWriter, r *http.Request) {
	var succ bool
	var data graph.ElementProperty
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		system.Logln("Failed to decode json")
	}
	vertexid, e := uuid.NewUUID()
	common.NoError(e)

	//partitions, err := server.Metadata.GetAllPartitions()
	//if err != nil {
	//	system.Logln("Failed to get partitions from zookeeper")
	//	fmt.Fprintf(w, "Error! Failed to add vertex")
	//}
	count := 0
	for _, partition := range partitions {
		if count == 3 {
			break
		}
		backendInfo, err := server.Metadata.GetBackendInformation(backend)
		if err != nil {
			continue
		}
		backendAddr, ok := backendInfo["address"]
		if !ok {
			continue
		}
		stClient, err := server.getOrCreateStorageClient(backendAddr.(string))
		if err != nil {
			continue
		}
		stClient.StoreVertex(graph.V(vertexid, data), &succ)
		if succ {
			count += 1
		}
	}

}

func (server *Server) deletevertex(w http.ResponseWriter, r *http.Request) {
	// Delete all edges that have destination as this vertex
	// Then delete the source vertex
	panic("todo")
}

func (server *Server) addedge(w http.ResponseWriter, r *http.Request) {
	panic("todo")
}

func (server *Server) deleteedge(w http.ResponseWriter, r *http.Request) {
	panic("todo")
}

func (server *Server) addproperty(w http.ResponseWriter, r *http.Request) {
	panic("todo")
}

func (server *Server) getsrcvertex(w http.ResponseWriter, r *http.Request) {
	panic("todo")
}

func (server *Server) getdestvertex(w http.ResponseWriter, r *http.Request) {
	panic("todo")
}

func (server *Server) getinedges(w http.ResponseWriter, r *http.Request) {
	panic("todo")
}

func (server *Server) getoutedges(w http.ResponseWriter, r *http.Request) {
	panic("todo")
}

func (server *Server) getparentvertices(w http.ResponseWriter, r *http.Request) {
	panic("todo")
}

func (server *Server) getchildvertices(w http.ResponseWriter, r *http.Request) {
	panic("todo")
}

func (server *Server) Serve() error {
	//panic("todo")
	http.HandleFunc("/AddVertex", server.addvertex)
	http.HandleFunc("/DeleteVertex", server.deletevertex)
	http.HandleFunc("/AddEdge", server.addedge)
	http.HandleFunc("/DeleteEdge", server.deleteedge)
	http.HandleFunc("/AddProperty", server.addproperty)
	http.HandleFunc("/GetSrcVertex", server.getsrcvertex)
	http.HandleFunc("/GetDestVertex", server.getdestvertex)
	http.HandleFunc("/GetInEdges", server.getinedges)
	http.HandleFunc("/GetOutEdges", server.getoutedges)
	http.HandleFunc("/GetParentVertices", server.getparentvertices)
	http.HandleFunc("/GetChildVertices", server.getchildvertices)

	go func() {
		log.Fatal(http.ListenAndServe(server.Config.Addr, nil))
	}()
	return nil
}
