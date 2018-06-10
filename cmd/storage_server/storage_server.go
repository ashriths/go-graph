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
)

var (
	storageAddr = flag.String("addr", "localhost:rand", "storage_server listen address")
	store       = flag.String("store", "memory", "storage_server data location")
	frc         = flag.String("rc", cmd.DefaultRCPath, "config file")
)

func GetStore(storeType string, metadataAddrs []string) storage.IOMapper {
	switch storeType {
	case "memory":
		return storage.NewInMemoryIOMapper(metadataAddrs)
	}
	panic("Invalid StoreType")
}

func main() {
	flag.Parse()
	args := flag.Args()
	system.Logging = true
	n := 0
	if len(args) == 0 {
		rc, e := cmd.LoadRC(*frc)
		common.NoError(e)
		run := func(i int) {
			if i > len(rc.Storage) {
				common.NoError(fmt.Errorf("back-end index out of range: %d", i))
			}

			backConfig := rc.StorageConfig(i, GetStore(*store, rc.MetadataServers))

			log.Printf("bin storage_server back-end serving on %s", backConfig.Addr)
			common.NoError(storage.ServeStorage(backConfig))
		}

		for i, b := range rc.Storage {

			if local.Check(b) {
				log.Println(i, b)
				go run(i)
				n++
			}
		}
	}

	if n > 0 {
		select {}
	}
}
