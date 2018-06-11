package main

import (
	"flag"
	"fmt"
	"github.com/ashriths/go-graph/cmd"
	"github.com/ashriths/go-graph/common"
	"github.com/ashriths/go-graph/local"
	"github.com/ashriths/go-graph/server"
	"github.com/ashriths/go-graph/system"
	"log"
)

var (
	serverAddr = flag.String("addr", "localhost:rand", "storage_server listen address")
	frc        = flag.String("rc", cmd.DefaultRCPath, "config file")
)

func main() {
	flag.Parse()
	args := flag.Args()
	system.Logging = true
	n := 0
	if len(args) == 0 {
		rc, e := cmd.LoadRC(*frc)
		common.NoError(e)
		run := func(i int) {
			if i > len(rc.Servers) {
				common.NoError(fmt.Errorf("back-end index out of range: %d", i))
			}
			serverConfig := rc.ServerConfig(i)

			log.Printf("server serving on %s", serverConfig.Addr)
			_server := server.NewZookeeperServer(serverConfig)
			common.NoError(_server.Serve())
		}

		for i, b := range rc.Servers {

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
