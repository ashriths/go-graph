package server

import (
	"log"
	"net/http"
	"reflect"
	"github.com/ashriths/go-graph/metadata"
	"github.com/ashriths/go-graph/storage"
)

type Server struct {
	Config  	  *ServerConfig
	ZkConnMap	  *metadata.ZkMetadataMapper
	storageClients map[string]*storage.StorageClient
}

type ServerConfig struct {
	MetadataServers []string
	Addr            string
	Ready           chan<- bool
}

func NewServer(config *ServerConfig) (error, *Server) {
	zkConnMap := &metadata.ZkMetadataMapper{ZkAddrs: config.MetadataServers}
	return nil, &Server{Config: config, ZkConnMap: zkConnMap}
}

func (server *Server) ZkCall(method string, args ...interface{}) ([]interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	out := reflect.ValueOf(server.ZkConnMap).MethodByName(method).Call(inputs)

	var output = make([]interface{}, len(out))
	for _,outp := range out[:len(out)] {
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

func (server *Server) addvertex() http.Handler {
	panic("todo")
}

func (server *Server) deletevertex() http.Handler {
	panic("todo")
}

func (server *Server) addedge() http.Handler {
	panic("todo")
}

func (server *Server) deleteedge() http.Handler {
	panic("todo")
}

func (server *Server) addproperty() http.Handler {
	panic("todo")
}

func (server *Server) getsrcvertex() http.Handler {
	panic("todo")
}

func (server *Server) getdestvertex() http.Handler {
	panic("todo")
}

func (server *Server) getinedges() http.Handler {
	panic("todo")
}

func (server *Server) getoutedges() http.Handler {
	panic("todo")
}

func (server *Server) getparentvertices() http.Handler {
	panic("todo")
}

func (server *Server) getchildvertices() http.Handler {
	panic("todo")
}

func (server *Server) Serve() error {
	//panic("todo")
	http.Handle("/AddVertex", server.addvertex())
	http.Handle("/DeleteVertex", server.deletevertex())
	http.Handle("/AddEdge", server.addedge())
	http.Handle("/DeleteEdge", server.deleteedge())
	http.Handle("/AddProperty", server.addproperty())
	http.Handle("/GetSrcVertex", server.getsrcvertex())
	http.Handle("/GetDestVertex", server.getdestvertex())
	http.Handle("/GetInEdges", server.getinedges())
	http.Handle("/GetOutEdges", server.getoutedges())
	http.Handle("/GetParentVertices", server.getparentvertices())
	http.Handle("/GetChildVertices", server.getchildvertices())

	go func() {
		log.Fatal(http.ListenAndServe(server.Config.Addr, nil))
	}()
	return nil
}
