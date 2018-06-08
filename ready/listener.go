package ready

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type listener struct {
	C chan string
}

func newListener() *listener {
	ret := new(listener)
	ret.C = make(chan string, 10)

	return ret
}

func (self *listener) Ready(s string, b *bool) error {
	self.C <- s
	*b = true

	return nil
}

func Listen(addr string) (<-chan string, error) {
	conn, e := net.Listen("tcp", addr)
	if e != nil {
		return nil, e
	}

	s := rpc.NewServer()
	lis := newListener()
	e = s.RegisterName("Ready", lis)
	if e != nil {
		return nil, e
	}

	go func() {
		e := http.Serve(conn, s)
		if e != nil {
			log.Print(e)
		}
	}()

	return lis.C, nil
}
