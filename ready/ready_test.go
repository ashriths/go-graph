package ready

import (
	"testing"
)

func TestReady(t *testing.T) {
	addr := "localhost:3000"

	ne := func(e error) {
		if e != nil {
			t.Fatal(e)
		}
	}

	as := func(cond bool) {
		if !cond {
			t.Fatal("assertion failed")
		}
	}

	c, e := Listen(addr)
	ne(e)

	msg := "So long and thanks for all the fish"
	ne(Notify(addr, msg))

	as(len(c) == 1)
	s := <-c

	as(s == msg)

	csend := Chan(addr, msg)

	csend <- true
	s = <-c
	as(s == msg)

	csend <- false
	s = <-c
	as(s == "!"+msg)
}
