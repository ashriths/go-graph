package main

import (
	"flag"
	"github.com/ashriths/go-graph/system"
)

var (
	serverAddr = flag.String("addr", "localhost:rand", "server listen address")
)

func main() {
	flag.Parse()

	*serverAddr = system.Resolve(*serverAddr)
	panic("Todo")

}
