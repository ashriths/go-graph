package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/locator"
	"github.com/ashriths/go-graph/metadata"
	"github.com/ashriths/go-graph/query"
	"github.com/ashriths/go-graph/storage"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
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

func NewZookeeperServer(config *ServerConfig) *Server {
	zkConnMap := &metadata.ZkMetadataMapper{ZkAddrs: config.MetadataServers}
	locator := &locator.RandomLocator{Metadata: zkConnMap}
	hqp := query.NewHTTPQueryParser()
	return &Server{
		Config:         config,
		Metadata:       zkConnMap,
		Locator:        locator,
		Parser:         hqp,
		storageClients: make(map[string]*storage.StorageClient),
	}
}

func (server *Server) ZkCall(method string, args ...interface{}) []interface{} {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	out := reflect.ValueOf(server.Metadata).MethodByName(method).Call(inputs)

	var output []interface{}
	for _, outp := range out {
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
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(byteData))
}

func handleError(w http.ResponseWriter, msg string) {
	system.Logln(msg)
	data := map[string]interface{}{"error": msg, "success": false}
	w.Header().Set("Content-Type", "application/json")
	writeResponse(w, data)
}

func (server *Server) addGraph(w http.ResponseWriter, r *http.Request) {
	graphId, e := uuid.NewUUID()
	if e != nil {
		handleError(w, "Cannot create UUID")
		return
	}
	if e := server.Metadata.CreateGraph(graphId, "{}"); e != nil {
		handleError(w, fmt.Sprintf("Failed to create graph. %s", e))
		return
	}
	data := map[string]interface{}{"data": map[string]string{"id": graphId.String()}, "success": true}
	writeResponse(w, data)
}

func (server *Server) addVertex(w http.ResponseWriter, r *http.Request) {
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

	partitionID, err := server.Locator.FindPartition(vertex)
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
			handleError(w, "Failed to create a storage client to backend: "+backendAddr)
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

func (server *Server) deleteVertex(w http.ResponseWriter, r *http.Request) {
	// Delete all edges that have destination as this vertex
	// Then delete the source vertex
	var succ bool
	var graphID, vertexID uuid.UUID

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

	vertexID_str, err := server.Parser.RetrieveParamByName(r, "vertexid")
	if err != nil {
		handleError(w, "Failed to fetch vertexid from request")
		return
	}
	vertexID, err = uuid.Parse(vertexID_str)
	if err != nil {
		handleError(w, "Failed to parse vertexid")
		return
	}

	edgeIds, err := server.Metadata.GetVertexInformation(graphID, vertexID)
	if err != nil {
		handleError(w, "Failed to get edge ids for vertex")
		return
	}

	for edgeIdStr, backends := range edgeIds {
		for _, backend := range backends.([]string) {
			stClient, err := server.getOrCreateStorageClient(backend)
			if err != nil {
				handleError(w, "Failed to get storage client")
				return
			}
			edgeId, err := uuid.Parse(edgeIdStr)
			if err != nil {
				handleError(w, "Failed to parse uuid")
				return
			}
			err = stClient.RemoveEdge(edgeId, &succ)
			if err != nil || succ == false {
				handleError(w, "Failed to remove inedges")
				return
			}
		}
	}

	_, backends, err := server.Metadata.GetVertexLocation(graphID, vertexID)
	if err != nil {
		handleError(w, "Failed to get backend infor for vertex")
		return
	}
	for _, backend := range backends {
		stClient, err := server.getOrCreateStorageClient(backend)
		if err != nil {
			handleError(w, "Failed to get storage client")
			return
		}
		err = stClient.RemoveVertex(vertexID, &succ)
		if err != nil || succ == false {
			handleError(w, "Failed to remove vertex")
			return
		}
	}
	err = server.Metadata.DeleteVertex(graphID, vertexID)
	if err != nil {
		handleError(w, "Failed to remove vertex from zookeeper")
		return
	}

	system.Logln("Successfully deleted vertex")

	responsemsg := map[string]interface{}{"msg": "Successfully deleted vertex", "success": true}
	writeResponse(w, responsemsg)
}

func (server *Server) addEdge(w http.ResponseWriter, r *http.Request) {
	var succ bool
	var data graph.ElementProperty
	var graphID uuid.UUID

	graphIdStr, err := server.Parser.RetrieveParamByName(r, "graphid")
	if err != nil {
		handleError(w, "Failed to fetch graphid from request")
		return
	}
	graphID, err = uuid.Parse(graphIdStr)
	if err != nil {
		handleError(w, "Failed to parse graphid")
		return
	}

	data, err = server.Parser.RetreiveQueryData(r)
	if err != nil {
		handleError(w, "Failed to parse request body")
		return
	}
	srcId, err := uuid.Parse(data["SrcVertex"])
	if err != nil {
		handleError(w, "Invalid UUID for SrcVertex")
		return
	}
	destId, err := uuid.Parse(data["DestVertex"])
	if err != nil {
		handleError(w, "Invalid UUID for DestVertex")
		return
	}
	edgeName, ok := data["Name"]
	if !ok {
		handleError(w, "Edge must have a relation Name")
		return
	}
	edgeId := uuid.New()
	edge := graph.E(graphID, edgeId, srcId, destId, edgeName, data)

	_, backends, err := server.Metadata.GetVertexLocation(graphID, srcId)
	for _, backend := range backends {
		data, err := server.Metadata.GetBackendInformation(backend)
		if err != nil {
			handleError(w, "Failed to get backend Info")
			return
		}
		backendAddr := data["address"].(string)

		stClient, err := server.getOrCreateStorageClient(backendAddr)
		if err != nil {
			handleError(w, "Failed to create a storage client to backend: "+backendAddr)
			return
		}

		e := stClient.StoreEdge(edge, &succ)
		if e != nil || !succ {
			handleError(w, "Failed to add edge to backend")
			return
		}
	}

	err = server.Metadata.CreateEdge(graphID, edgeId, srcId)
	if err != nil {
		handleError(w, "Failed to create edge in Metadata")
		return
	}

	//var inedge = map[string][]string{}
	_, backends, err = server.Metadata.GetVertexLocation(graphID, destId)
	if err != nil {
		handleError(w, "Failed to fetch backends for destination vertex")
		return
	}
	var backendInfo []string
	for _, backend := range backends {
		data, err := server.Metadata.GetBackendInformation(backend)
		if err != nil {
			handleError(w, "Failed to fetch backend info for backend")
			return
		}
		backendInfo = append(backendInfo, data["address"].(string))
	}

	//inedge[edgeId.String()] = backendInfo
	err = server.Metadata.UpdateVertexInformation(graphID, destId, edgeId.String(), backendInfo)
	if err != nil {
		handleError(w, "Failed to add vertex information")
		return
	}

	err = server.Metadata.UpdateEdgeInformation(graphID, edgeId, "dest", destId.String())
	if err != nil {
		handleError(w, "Failed to add edge information")
		return
	}

	system.Logln("Successfully added Edge")

	responsedata := map[string]interface{}{"edgeId": edgeId.String()}
	responsemsg := map[string]interface{}{"msg": "Successfully added edge", "success": true, "data": responsedata}
	writeResponse(w, responsemsg)

}

func (server *Server) deleteEdge(w http.ResponseWriter, r *http.Request) {
	//panic("todo")
	var succ bool
	var graphID, edgeID uuid.UUID

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

	edgeID_str, err := server.Parser.RetrieveParamByName(r, "edgeid")
	if err != nil {
		handleError(w, "Failed to fetch edgeid from request")
		return
	}
	edgeID, err = uuid.Parse(edgeID_str)
	if err != nil {
		handleError(w, "Failed to parse edgeid")
		return
	}

	server.findAndRunRPCOnBackend(w, graphID, edgeID, "Edge", "RemoveEdge", edgeID, &succ)
	err = server.Metadata.DeleteEdge(graphID, edgeID)
	if err != nil {
		handleError(w, "Failed to delete edge znode on zookeeper")
		return
	}

	destIdStr, err := server.Metadata.GetEdgeInformation(graphID, edgeID)
	if err != nil {
		handleError(w, "Failed to get edge information")
		return
	}

	destId, err := uuid.Parse(destIdStr["dest"].(string))
	if err != nil {
		handleError(w, "Failed to parse uuid")
		return
	}

	err = server.Metadata.DeleteVertexInformation(graphID, destId, "dest")
	if err != nil {
		handleError(w, "Failed to add vertex information")
		return
	}

	system.Logln("Successfully deleted edge")

	responsemsg := map[string]interface{}{"msg": "Successfully deleted edge", "success": true}
	writeResponse(w, responsemsg)
}

func (server *Server) getInEdges(w http.ResponseWriter, r *http.Request) {
	//panic("todo")
	var edges []graph.Edge
	var graphID, vertexId uuid.UUID

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

	vertexIdStr, err := server.Parser.RetrieveParamByName(r, "vertexid")
	if err != nil {
		handleError(w, "Failed to fetch vertexid from request")
		return
	}
	vertexId, err = uuid.Parse(vertexIdStr)
	if err != nil {
		handleError(w, "Failed to parse vertexid")
		return
	}

	inEdges, err := server.Metadata.GetVertexInformation(graphID, vertexId)
	if err != nil {
		handleError(w, "Failed to get InEdgeIDs from zookeeper")
		return
	}

	for _, backends := range inEdges {
		for _, backend := range backends.([]string) {
			data, err := server.Metadata.GetBackendInformation(backend)
			if err != nil {
				continue
			}
			backendAddr := data["address"].(string)

			stClient, err := server.getOrCreateStorageClient(backendAddr)
			if err != nil {
				continue
			}
			e := stClient.GetOutEdges(vertexId, &edges)
			if e != nil {
				continue
			}
			break
		}
	}
	if edges == nil {
		handleError(w, "Failed to fetch edgeid from request")
		return
	}
	system.Logln("Successfully got in edges")

	responsemsg := map[string]interface{}{"msg": "Successfully got in edges", "success": true, "data": edges}
	writeResponse(w, responsemsg)
}

func (server *Server) getOutEdges(w http.ResponseWriter, r *http.Request) {
	var edges []graph.Edge
	var graphID, vertexId uuid.UUID

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

	vertexIdStr, err := server.Parser.RetrieveParamByName(r, "vertexid")
	if err != nil {
		handleError(w, "Failed to fetch vertexid from request")
		return
	}
	vertexId, err = uuid.Parse(vertexIdStr)
	if err != nil {
		handleError(w, "Failed to parse vertexid")
		return
	}
	_, backends, err := server.Metadata.GetVertexLocation(graphID, vertexId)
	for _, backend := range backends {
		data, err := server.Metadata.GetBackendInformation(backend)
		if err != nil {
			continue
		}
		backendAddr := data["address"].(string)

		stClient, err := server.getOrCreateStorageClient(backendAddr)
		if err != nil {
			continue
		}
		e := stClient.GetOutEdges(vertexId, &edges)
		if e != nil {
			continue
		}
	}
	if edges == nil {
		handleError(w, "Failed to fetch edgeid from request")
		return
	}
	system.Logln("Successfully got out edges")

	responsemsg := map[string]interface{}{"msg": "Successfully got out edges", "success": true, "data": edges}
	writeResponse(w, responsemsg)

}

func (server *Server) getParentVertices(w http.ResponseWriter, r *http.Request) {
	var edges []graph.Edge
	var vertices []graph.Vertex
	var graphID, vertexId uuid.UUID

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

	vertexIdStr, err := server.Parser.RetrieveParamByName(r, "vertexid")
	if err != nil {
		handleError(w, "Failed to fetch vertexid from request")
		return
	}
	vertexId, err = uuid.Parse(vertexIdStr)
	if err != nil {
		handleError(w, "Failed to parse vertexid")
		return
	}

	inEdges, err := server.Metadata.GetVertexInformation(graphID, vertexId)
	if err != nil {
		handleError(w, "Failed to get InEdgeIDs from zookeeper")
		return
	}

	for _, backends := range inEdges {
		for _, backend := range backends.([]string) {
			data, err := server.Metadata.GetBackendInformation(backend)
			if err != nil {
				continue
			}
			backendAddr := data["address"].(string)

			stClient, err := server.getOrCreateStorageClient(backendAddr)
			if err != nil {
				continue
			}
			e := stClient.GetOutEdges(vertexId, &edges)
			if e != nil {
				continue
			}
			break
		}
	}
	if edges == nil {
		handleError(w, "Failed to fetch edgeid from request")
		return
	}

	for _, edge := range edges {
		err, srcvertex := edge.GetSrcVertex()
		if err != nil {
			handleError(w, "Failed to get destination vertex of edge")
			return
		}

		var e error
		_, backends, err := server.Metadata.GetVertexLocation(graphID, srcvertex.GetUUID())
		for _, backend := range backends {
			var vertex graph.Vertex
			data, err := server.Metadata.GetBackendInformation(backend)
			if err != nil {
				continue
			}
			backendAddr := data["address"].(string)

			stClient, err := server.getOrCreateStorageClient(backendAddr)
			if err != nil {
				continue
			}
			e = stClient.GetVertexById(srcvertex.GetUUID(), &vertex)
			if e != nil {
				continue
			} else {
				vertices = append(vertices, vertex)
				break
			}
		}
		if e != nil {
			handleError(w, "Failed to get some parent vertices")
			return
		}
	}

	system.Logln("Successfully got parent vertices")

	responsemsg := map[string]interface{}{"msg": "Successfully got parent vertices", "success": true, "data": edges}
	writeResponse(w, responsemsg)
}

func (server *Server) getChildVertices(w http.ResponseWriter, r *http.Request) {
	var edges []graph.Edge
	var vertices []graph.Vertex
	var graphID, vertexId uuid.UUID

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

	vertexIdStr, err := server.Parser.RetrieveParamByName(r, "vertexid")
	if err != nil {
		handleError(w, "Failed to fetch vertexid from request")
		return
	}
	vertexId, err = uuid.Parse(vertexIdStr)
	if err != nil {
		handleError(w, "Failed to parse vertexid")
		return
	}
	_, backends, err := server.Metadata.GetVertexLocation(graphID, vertexId)
	for _, backend := range backends {
		data, err := server.Metadata.GetBackendInformation(backend)
		if err != nil {
			continue
		}
		backendAddr := data["address"].(string)

		stClient, err := server.getOrCreateStorageClient(backendAddr)
		if err != nil {
			continue
		}
		e := stClient.GetOutEdges(vertexId, &edges)
		if e != nil {
			continue
		}
	}
	if edges == nil {
		handleError(w, "Failed to fetch edgeid from request")
		return
	}

	for _, edge := range edges {
		err, destvertex := edge.GetDestVertex()
		if err != nil {
			handleError(w, "Failed to get destination vertex of edge")
			return
		}

		var e error
		_, backends, err := server.Metadata.GetVertexLocation(graphID, destvertex.GetUUID())
		for _, backend := range backends {
			var vertex graph.Vertex
			data, err := server.Metadata.GetBackendInformation(backend)
			if err != nil {
				continue
			}
			backendAddr := data["address"].(string)

			stClient, err := server.getOrCreateStorageClient(backendAddr)
			if err != nil {
				continue
			}
			e = stClient.GetVertexById(destvertex.GetUUID(), &vertex)
			if e != nil {
				continue
			} else {
				vertices = append(vertices, vertex)
				break
			}
		}
		if e != nil {
			handleError(w, "Failed to get some child vertices")
			return
		}
	}

	system.Logln("Successfully got child vertices")

	responsemsg := map[string]interface{}{"msg": "Successfully got child vertices", "success": true, "data": vertices}
	writeResponse(w, responsemsg)

}

func (server *Server) findAndRunRPCOnBackend(w http.ResponseWriter, graphID uuid.UUID,
	elementID uuid.UUID, elementType string, method string, args ...interface{}) {
	var err error
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	output := server.ZkCall("Get"+elementType+"Location", graphID, elementID)
	if output[2] != nil {
		err := output[2].(error)
		handleError(w, fmt.Sprintf("Failed to get location. %s", err))
	}
	backendids := output[1].([]string)
	for _, backendid := range backendids {
		backendInfo, err := server.Metadata.GetBackendInformation(backendid)
		backendAddr := backendInfo["address"].(string)
		stClient, err := server.getOrCreateStorageClient(backendAddr)
		if err != nil {
			continue
		}
		ret := reflect.ValueOf(stClient).MethodByName(method).Call(inputs)
		if ret[0].Interface() == nil {
			break
		}
		err = ret[0].Interface().(error)
	}
	if err != nil {
		handleError(w, "Failed to run query on backend")
	}
}

func (server *Server) getVertex(w http.ResponseWriter, r *http.Request) {
	var graphID uuid.UUID
	var vertexID uuid.UUID
	var vertex graph.Vertex

	graphID_str, err := server.Parser.RetrieveParamByName(r, "graphid")
	if err != nil {
		handleError(w, "Failed to fetch graphid from request")
		return
	}
	graphID, err = uuid.Parse(graphID_str)
	if err != nil {
		handleError(w, "Failed to parse graphid: "+graphID_str)
		return
	}

	vertexID_str, err := server.Parser.RetrieveParamByName(r, "vertexid")
	if err != nil {
		handleError(w, "Failed to fetch vertexid from request")
		return
	}
	vertexID, err = uuid.Parse(vertexID_str)
	if err != nil {
		handleError(w, "Failed to parse vertexid: "+vertexID_str)
		return
	}

	server.findAndRunRPCOnBackend(w, graphID, vertexID, "Vertex", "GetVertexById", vertexID, &vertex)
	responsemsg := map[string]interface{}{"msg": "Successfully fetches vertex properties", "success": true, "data": vertex}
	writeResponse(w, responsemsg)
}

func (server *Server) getEdge(w http.ResponseWriter, r *http.Request) {
	var graphID uuid.UUID
	var edgeID uuid.UUID
	var edge graph.Edge

	graphID_str, err := server.Parser.RetrieveParamByName(r, "graphid")
	if err != nil {
		handleError(w, "Failed to fetch graphid from request")
		return
	}
	graphID, err = uuid.Parse(graphID_str)
	if err != nil {
		handleError(w, "Failed to parse graphid: "+graphID_str)
		return
	}

	edgeID_str, err := server.Parser.RetrieveParamByName(r, "edgeid")
	if err != nil {
		handleError(w, "Failed to fetch edgeid from request")
		return
	}
	edgeID, err = uuid.Parse(edgeID_str)
	if err != nil {
		handleError(w, "Failed to parse edgeid: "+edgeID_str)
		return
	}

	server.findAndRunRPCOnBackend(w, graphID, edgeID, "Edge", "GetEdgeById", edgeID, &edge)

	responsemsg := map[string]interface{}{"msg": "Successfully fetches edge properties", "success": true, "data": edge}
	writeResponse(w, responsemsg)
}

func (server *Server) getGraphs(w http.ResponseWriter, r *http.Request) {
	var graphIds []uuid.UUID
	if e := server.Metadata.GetGraphs(&graphIds); e != nil {
		handleError(w, fmt.Sprintf("Failed to get graphs. %s", e))
		return
	}
	responsemsg := map[string]interface{}{"msg": "Successfully fetched graphs", "success": true, "data": graphIds}
	writeResponse(w, responsemsg)
}

func (server *Server) Serve() error {
	//panic("todo")
	http.HandleFunc("/AddGraph", server.addGraph)
	http.HandleFunc("/AddVertex", server.addVertex)
	http.HandleFunc("/DeleteVertex", server.deleteVertex)
	http.HandleFunc("/AddEdge", server.addEdge)
	http.HandleFunc("/DeleteEdge", server.deleteEdge)
	//http.HandleFunc("/updateVertex", server.updateVertex)
	//http.HandleFunc("/updateEdge", server.updateEdge)
	http.HandleFunc("/GetInEdges", server.getInEdges)
	http.HandleFunc("/GetOutEdges", server.getOutEdges)
	http.HandleFunc("/GetParentVertices", server.getParentVertices)
	http.HandleFunc("/GetChildVertices", server.getChildVertices)
	http.HandleFunc("/GetVertex", server.getVertex)
	http.HandleFunc("/GetEdge", server.getEdge)
	http.HandleFunc("/GetGraphs", server.getGraphs)

	go func() {
		log.Fatal(http.ListenAndServe(server.Config.Addr, nil))
	}()
	return nil
}

//var sc = &ServerConfig{
//	MetadataServers: []string{"169.228.66.172:21810", "169.228.66.170:21810", "169.228.66.171:21810"},
//	Addr: "0.0.0.0:12345",
//	Ready: make(chan bool),
//}
//var _ Server = NewZookeeperServer(sc)
