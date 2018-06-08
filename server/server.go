package server

import (
	"log"
	"net/http"
	"strings"
)

type Server struct {
	Config *ServerConfig
}

type ServerConfig struct {
	MetadataServers []string
	Addr            string
	Ready           chan<- bool
}

func NewServer(sc *ServerConfig) (error, *Server) {
	return nil, &Server{Addr: sc.Addr, Config: sc}
}

func (server *Server) addvertex() http.Handler {

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
	http.Handle("/AddVertex", addvertex())
	http.Handle("/DeleteVertex", deletevertex())
	http.Handle("/AddEdge", addedge())
	http.Handle("/DeleteEdge", deleteedge())
	http.Handle("/AddProperty", addproperty())
	http.Handle("/GetSrcVertex", getsrcvertex())
	http.Handle("/GetDestVertex", getdestvertex())
	http.Handle("/GetInEdges", getinedges())
	http.Handle("/GetOutEdges", getoutedges())
	http.Handle("/GetParentVertices", getparentvertices())
	http.Handle("/GetChildVertices", getchildvertices())

	go func() {
		log.Fatal(http.ListenAndServe(server.Config.Addr, nil))
	}()
}
