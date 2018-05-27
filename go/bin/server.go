package main

import (
	"flag"
	"go-graph/src/system"
	"go-graph/src/server"
)

var (
	addr = flag.String("addr", "localhost:rand", "server listen address")
)

func main() {
	flag.Parse()

	*addr = system.Resolve(*addr)
	panic("Todo")
	server.NewServer(&server.ServerConfig{
		Addr: addr,
		MetadataServers:
	})
}
