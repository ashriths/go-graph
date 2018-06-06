package rc

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var DefaultRCPath = "conf.rc"

type RC struct {
	Backs   []string
	Keepers []string
}

type BackAddr struct {
	Serve string
	Peer  string
}

func (self *RC) BackCount() int {
	return len(self.Backs)
}

func (self *RC) BackConfig(i int, s Storage) *BackConfig {
	ret := new(BackConfig)
	ret.Addr = self.Backs[i]
	ret.Store = s
	ret.Ready = make(chan bool, 1)

	return ret
}

func (self *RC) KeeperConfig(i int) *KeeperConfig {
	if i >= len(self.Keepers) {
		panic("keeper index out of range")
	}

	ret := new(KeeperConfig)
	ret.Backs = self.Backs
	ret.Addrs = self.Keepers
	ret.This = i
	ret.Id = time.Now().UnixNano()

	return ret
}

func LoadRC(p string) (*RC, error) {
	fin, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer fin.Close()

	ret := new(RC)
	e = json.NewDecoder(fin).Decode(ret)
	if e != nil {
		return nil, e
	}

	return ret, nil
}

func (self *RC) marshal() []byte {
	b, e := json.MarshalIndent(self, "", "    ")
	if e != nil {
		panic(e)
	}

	return b
}

func (self *RC) Save(p string) error {
	b := self.marshal()

	fout, e := os.Create(p)
	if e != nil {
		return e
	}

	_, e = fout.Write(b)
	if e != nil {
		return e
	}

	_, e = fmt.Fprintln(fout)
	if e != nil {
		return e
	}

	return fout.Close()
}

func (self *RC) String() string {
	b := self.marshal()
	return string(b)
}
