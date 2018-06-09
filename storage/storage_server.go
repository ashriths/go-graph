package storage

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func ServeStorage(storageConfig *StorageConfig) error {
	rpcServer := rpc.NewServer()
	e := rpcServer.RegisterName("Storage", storageConfig.Store)
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
