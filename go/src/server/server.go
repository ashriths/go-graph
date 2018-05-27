package server

type Server struct {
	Addr string
	Config *ServerConfig
}

type ServerConfig struct {
	MetadataServers []string
	Addr string
	Ready chan<- bool
}

func NewServer(sc *ServerConfig)  (error, *Server) {
	return nil,&Server{Addr:sc.Addr, Config:sc}
}

func (server *Server)Serve() error {
	panic("todo")
}