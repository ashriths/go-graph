package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/ashriths/go-graph/storage"
)

var DefaultRCPath = "conf.rc"

type RC struct {
	Storage   []string
	MetadataServers []string
}

type BackAddr struct {
	Serve string
	Peer  string
}

func (self *RC) StorageCount() int {
	return len(self.Storage)
}

func (self *RC) StorageConfig(i int, s storage.IOMapper) *storage.StorageConfig {
	ret := new(storage.StorageConfig)
	ret.Addr = self.Storage[i]
	ret.Store = s
	ret.Ready = make(chan bool, 1)
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
