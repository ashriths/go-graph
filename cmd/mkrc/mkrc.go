package mkrc

import (
	"flag"
	"fmt"
	"log"
	"go-graph/system"
	"github.com/ashriths/go-graph/cmd"
)

// For now, we assume that we have sequentially-IP'd hosts that don't span more
// than one octet.
const IP_PREFIX = "169.228.66"
const FIRST_IP = 166
const NUM_HOSTS = 10

var (
	local   = flag.Bool("local", false, "always use local ports")
	nserver   = flag.Int("nserver", 1, "number of back ends servers")
	nmetadata   = flag.Int("nmetadata", 1, "number of metadata backends")
	full    = flag.Bool("full", false, "setup of 10 back-ends and 3 keepers")
	fixPort = flag.Bool("fix", false, "fix port numbers; don't use random ones")
	file     = flag.String("file", cmd.DefaultRCPath, "config file path")
)

func main() {
	flag.Parse()

	if *nserver > 300 {
		log.Fatal(fmt.Errorf("too many back-ends"))
	}
	if *nmetadata > NUM_HOSTS {
		log.Fatal(fmt.Errorf("too many metadata back-ends"))
	}

	if *full {
		*nserver = NUM_HOSTS
		*nmetadata = 3
	}

	p := 3000
	if !*fixPort {
		p = system.RandPort()
	}

	rc := new(cmd.RC)
	rc.Storage = make([]string, *nserver)
	rc.MetadataServers = make([]string, *nmetadata)

	if !*local {
		const ipOffset = FIRST_IP
		const nmachine = NUM_HOSTS

		for i := 0; i < *nserver; i++ {
			host := fmt.Sprintf("%s.%d", IP_PREFIX, ipOffset+i%nmachine)
			rc.Storage[i] = fmt.Sprintf("%s:%d", host, p+i/nmachine)
		}

		p += *nserver / nmachine
		if *nserver%nmachine > 0 {
			p++
		}

		for i := 0; i < *nmetadata; i++ {
			host := fmt.Sprintf("%s.%d", IP_PREFIX, ipOffset+i%nmachine)
			rc.MetadataServers[i] = fmt.Sprintf("%s:%d", host, p)
		}
	} else {
		for i := 0; i < *nserver; i++ {
			rc.Storage[i] = fmt.Sprintf("localhost:%d", p)
			p++
		}

		for i := 0; i < *nmetadata; i++ {
			rc.MetadataServers[i] = fmt.Sprintf("localhost:%d", p)
			p++
		}
	}

	fmt.Println(rc.String())

	if *file != "" {
		e := rc.Save(*file)
		if e != nil {
			log.Fatal(e)
		}
	}
}
