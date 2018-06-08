package ready

import (
	"log"
	"net/rpc"
)

func Notify(addr, s string) error {
	c, e := rpc.DialHTTP("tcp", addr)
	if e != nil {
		return e
	}

	var b bool
	e = c.Call("Ready.Ready", s, &b)
	if e != nil {
		return e
	}

	return c.Close()
}

func NotifyFail(addr, s string) error {
	return Notify(addr, "!"+s)
}

func notify(c chan bool, addr, s string) {
	var e error
	for {
		b := <-c
		if b {
			e = Notify(addr, s)
		} else {
			e = NotifyFail(addr, s)
		}

		if e != nil {
			log.Print(e)
		}
	}
}

func Chan(addr, s string) chan<- bool {
	c := make(chan bool)
	go notify(c, addr, s)
	return c
}
