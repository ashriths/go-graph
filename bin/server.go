package bin

import (
	"flag"
)

var (
	addr = flag.String("addr", "localhost:rand", "server listen address")
)

func main() {
	flag.Parse()

	*addr = randaddr.Resolve(*addr)
}
