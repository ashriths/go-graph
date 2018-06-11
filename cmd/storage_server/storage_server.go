package main

import (
	"flag"
	"fmt"
	"github.com/ashriths/go-graph/cmd"
	"github.com/ashriths/go-graph/common"
	"github.com/ashriths/go-graph/local"
	"github.com/ashriths/go-graph/storage"
	"github.com/ashriths/go-graph/system"
	"log"
	"strconv"
)

var (
	storageAddr = flag.String("addr", "localhost:rand", "storage_server listen address")
	store       = flag.String("store", "memory", "storage_server data location")
	frc         = flag.String("rc", cmd.DefaultRCPath, "config file")
)

func GetStore(storeType string, backendId string, metadataAddrs []string) storage.IOMapper {
	switch storeType {
	case "memory":
		return storage.NewInMemoryIOMapper(backendId, metadataAddrs)
	}
	panic("Invalid StoreType")
}

func main() {
	flag.Parse()
	args := flag.Args()
	system.Logging = true
	n := 0

	rc, e := cmd.LoadRC(*frc)
	common.NoError(e)

	run := func(i int) {
		if i > len(rc.Storage) {
			common.NoError(fmt.Errorf("back-end index out of range: %d", i))
		}
		backendId, e := storage.RegisterStorage(rc.Storage[i], rc.MetadataServers)
		if e != nil {
			log.Fatal("Unable to register storage.")
		}
		storageConfig := rc.StorageConfig(i, GetStore(*store, backendId, rc.MetadataServers))

		log.Printf("bin storage_server back-end serving on %s", storageConfig.Addr)
		common.NoError(storage.ServeStorage(storageConfig))
	}

	if len(args) == 0 {

		for i, b := range rc.Storage {

			if local.Check(b) {
				log.Println(i, b)
				go run(i)
				n++
			}
		}
	} else {
		// scan for indices for the addresses
		for _, a := range args {
			i, e := strconv.Atoi(a)
			common.NoError(e)
			go run(i)
			n++
		}
	}

	if n > 0 {
		select {}
	}
}
