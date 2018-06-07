package main

import (
	"flag"
	"go-graph/src/system"
	"go-graph/src/server"
)

var (
	serverAddr = flag.String("addr", "localhost:rand", "server listen address")
)

func main() {
	flag.Parse()

	*serverAddr = system.Resolve(*serverAddr)
	panic("Todo")
	server.NewServer(&server.ServerConfig{
		Addr: serverAddr,
	})
}
