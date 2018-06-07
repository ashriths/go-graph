package cmd

import (
	"flag"
	"go-graph/go/src/system"
	"net/rpc"
	"fmt"
	"net"
	"net/http"
	"go-graph/go/src/storage"
)

var (
	storageAddr = flag.String("addr", "localhost:rand", "storage listen address")
	store = flag.String("store", "memory", "storage data location")
)

func ServeStorage(storageConfig *storage.StorageConfig) error{
	server := rpc.NewServer()
	e := server.Register(storageConfig.Store)
	if e != nil {
		fmt.Println(e)
		if storageConfig.Ready != nil {
			storageConfig.Ready <- false
		}
		return e
	}
	//rpc.HandleHTTP()

	l, e := net.Listen("tcp", storageConfig.Addr)
	if e != nil {
		fmt.Println(e)
		if storageConfig.Ready != nil {
			storageConfig.Ready <- false
		}
		return e
	}

	if storageConfig.Ready != nil {
		storageConfig.Ready <- true
	}

	return http.Serve(l, server)
}

func GetStore(storeType string)  storage.IOMapper{
	switch storeType {
	case "memory":
		return storage.NewInMemoryIOMapper()
	}
	panic("Invalid StoreType")
}

func main() {
	flag.Parse()
	address := system.Resolve(*storageAddr)
	conf := &storage.StorageConfig{
		Store: GetStore(*store),
		Addr: address,
		Ready: make(chan bool),
	}
	ServeStorage(conf)
}
