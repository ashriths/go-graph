package server

import (
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/locator"
	"github.com/ashriths/go-graph/metadata"
	"github.com/ashriths/go-graph/storage"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
	"log"
	"net/http"
	"reflect"
	"fmt"
	"github.com/ashriths/go-graph/query"
	"encoding/json"
)

type Server struct {
	Config         *ServerConfig
	Metadata       metadata.Metadata
	storageClients map[string]*storage.StorageClient
	Locator        locator.Locator
	Parser         query.QueryParser
}

type ServerConfig struct {
	MetadataServers []string
	Addr            string
	Ready           chan<- bool
}

const (
	LOCATOR = "RandomLocator"
	REPLICATIONFACTOR = 3
)

func NewZookeeperServer(config *ServerConfig) (error, *Server) {
	zkConnMap := &metadata.ZkMetadataMapper{ZkAddrs: config.MetadataServers}
	locator := &locator.RandomLocator{Metadata: zkConnMap}
	hqp := query.NewHTTPQueryParser()
	return nil, &Server{Config: config, Metadata: zkConnMap, Locator: locator, Parser: hqp}
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

func writeResponse(w http.ResponseWriter, data map[string]interface{}) {
	byteData, err := json.Marshal(data)
	if err != nil {
		system.Logln("Failed to marshal response data")
	}
	fmt.Fprintf(w, string(byteData))
}

func handleError(w http.ResponseWriter, msg string) {
	system.Logln(msg)
	data := map[string]interface{}{"error": msg, "success": false}
	writeResponse(w, data)
}

func (server *Server) addvertex(w http.ResponseWriter, r *http.Request) {
	var succ bool
	var data graph.ElementProperty
	var graphID uuid.UUID

	graphID_str, err := server.Parser.RetrieveParamByName(r, "graphid")
	if err != nil {
		handleError(w, "Failed to fetch graphid from request")
		return
	}
	graphID, err = uuid.Parse(graphID_str)
	if err != nil {
		handleError(w, "Failed to parse graphid")
		return
	}

	data, err = server.Parser.RetreiveQueryData(r)
	if err != nil {
		handleError(w, "Failed to parse request body")
		return
	}
	vertexid := uuid.New()
	vertex := graph.V(graphID, vertexid, data)

	partitionID, err := server.Locator.FindPartition(vertex.Element)
	if err != nil {
		handleError(w, "Failed to get a partition")
		return
	}

	backends, err := server.Metadata.GetBackendsForPartition(graphID, partitionID)
	for _, backend := range backends {
		data, err := server.Metadata.GetBackendInformation(backend)
		if err != nil {
			handleError(w, "Failed to get backend Info")
			return
		}
		backendAddr := data["address"].(string)

		stClient, err := server.getOrCreateStorageClient(backendAddr)
		if err != nil {
			handleError(w, "Failed to create a storage client to backend: " + backendAddr)
			return
		}

		e := stClient.StoreVertex(vertex, &succ)
		if e != nil || !succ {
			handleError(w, "Failed to add vertex to backend")
			return
		}
	}

	err = server.Metadata.CreateVertex(graphID, partitionID, vertexid)
	if err != nil {
		handleError(w, "Failed to create vertex zNode in zookeeper")
		return
	}

	system.Logln("Successfully added vertex")

	responsedata := map[string]interface{}{"vertexID": vertexid.String()}
	responsemsg := map[string]interface{}{"msg": "Successfully added vertex", "success": true, "data": responsedata}
	writeResponse(w, responsemsg)
}

func (server *Server) deletevertex(w http.ResponseWriter, r *http.Request) {
	// Delete all edges that have destination as this vertex
	// Then delete the source vertex
	panic("todo")
}

func (server *Server) addedge(w http.ResponseWriter, r *http.Request) {
	//panic("todo")

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

func (server *Server) getvertexproperties(w http.ResponseWriter, r *http.Request) {
	var graphID uuid.UUID
	var vertexID uuid.UUID
	var vertex *graph.Vertex

	graphID_str, err := server.Parser.RetrieveParamByName(r, "graphid")
	if err != nil {
		handleError(w, "Failed to fetch graphid from request")
		return
	}
	graphID, err = uuid.Parse(graphID_str)
	if err != nil {
		handleError(w, "Failed to parse graphid: " + graphID_str)
		return
	}

	vertexID_str, err := server.Parser.RetrieveParamByName(r, "vertexid")
	if err != nil {
		handleError(w, "Failed to fetch vertexid from request")
		return
	}
	vertexID, err = uuid.Parse(vertexID_str)
	if err != nil {
		handleError(w, "Failed to parse vertexid: " + vertexID_str)
		return
	}

	backendids, err := server.Metadata.GetVertexLocation(graphID, vertexID)
	for _, backendid := range backendids {
		 backendInfo, err := server.Metadata.GetBackendInformation(backendid)
		 backendAddr := backendInfo["address"].(string)
		 stClient, err :=  server.getOrCreateStorageClient(backendAddr)
		 if err != nil {
		 	continue
		 }
		 err = stClient.GetVertexById(vertexID, vertex)
		 if err == nil {
		 	break
		 }
	}

	err, properties := vertex.GetProperties()
	if err != nil {
		handleError(w, "Failed to fetch vertex properties")
		return
	}
	responsedata := map[string]interface{}{"VertexProperties": properties}
	responsemsg := map[string]interface{}{"msg": "Successfully fetches vertex properties", "success": true, "data": responsedata}
	writeResponse(w, responsemsg)
}

func (server *Server) getedgeproperties(w http.ResponseWriter, r *http.Request) {

}

func (server *Server) setproperties(w http.ResponseWriter, r *http.Request) {

}

func (server *Server) getgraphid(w http.ResponseWriter, r *http.Request) {

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
	http.HandleFunc("/GetVertexProperties", server.getvertexproperties)
	http.HandleFunc("/GetEdgeProperties", server.getedgeproperties)
	http.HandleFunc("/SetProperties", server.setproperties)
	http.HandleFunc("/GetGraphId", server.getgraphid)


	go func() {
		log.Fatal(http.ListenAndServe(server.Config.Addr, nil))
	}()
	return nil
}

var sc *ServerConfig = &ServerConfig{MetadataServers: []string{"169.228.66.172:21810", "169.228.66.170:21810", "169.228.66.171:21810"}, Addr: "0.0.0.0:12345", Ready: make(chan bool)}
var server, _ = NewZookeeperServer(sc)