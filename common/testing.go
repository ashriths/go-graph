package common

import (
	"runtime/debug"
	"testing"
)

func Assert(cond bool, t *testing.T) {
	if !cond {
		debug.PrintStack()
		t.Fatal("assertion failed")
	}
}
