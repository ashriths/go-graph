package storage

import (
	"fmt"
	"github.com/ashriths/go-graph/common"
	"github.com/ashriths/go-graph/metadata"
	"github.com/ashriths/go-graph/system"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func RegisterStorage(addr string, metadataAddrs []string) (string, error) {
	metadataConnection := metadata.NewZkMetadataMapper(metadataAddrs)
	backendId, e := metadataConnection.CreateBackend(addr)
	if e != nil {
		common.LogError(e)
		return "", e
	}
	system.Logln("Registered backend to ", metadataAddrs, " : ", backendId)
	return backendId, nil
}

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
