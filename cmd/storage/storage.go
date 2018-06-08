package main

import (
	"flag"
	"net/rpc"
	"fmt"
	"net"
	"net/http"
	"github.com/ashriths/go-graph/storage"
	"log"
	"github.com/ashriths/go-graph/common"
	"github.com/ashriths/go-graph/cmd"
	"github.com/ashriths/go-graph/local"
)

var (
	storageAddr = flag.String("addr", "localhost:rand", "storage listen address")
	store = flag.String("store", "memory", "storage data location")
	frc = flag.String("rc", cmd.DefaultRCPath, "config file")
)

func ServeStorage(storageConfig *storage.StorageConfig) error{
	rpcServer := rpc.NewServer()
	e := rpcServer.Register(storageConfig.Store)
	if e != nil {
		fmt.Println(e)
		if storageConfig.Ready != nil {
			storageConfig.Ready <- false
		}
		return e
	}
	//rpc.HandleHTTP()
	log.Println("Registered store")
	l, e := net.Listen("tcp", storageConfig.Addr)
	if e != nil {
		fmt.Println(e)
		if storageConfig.Ready != nil {
			storageConfig.Ready <- false
		}
		return e
	}
	log.Println("Listener configured")

	if storageConfig.Ready != nil {
		storageConfig.Ready <- true
	}
	log.Println("Storage RPC server started on ", storageConfig.Addr)
	return http.Serve(l, rpcServer)
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
	args := flag.Args()

	n := 0
	if len(args) == 0 {
		rc, e := cmd.LoadRC(*frc)
		common.NoError(e)
		run := func(i int) {
			if i > len(rc.Storage) {
				common.NoError(fmt.Errorf("back-end index out of range: %d", i))
			}

			backConfig := rc.StorageConfig(i, GetStore(*store))

			log.Printf("bin storage back-end serving on %s", backConfig.Addr)
			common.NoError(ServeStorage(backConfig))
		}

		for i, b := range rc.Storage {

			if local.Check(b) {
				log.Println(i,b)
				go run(i)
				n++
			}
		}
	}

	if n > 0 {
		select {}
	}
}
