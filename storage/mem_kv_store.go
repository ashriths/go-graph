package storage

import (
	"container/list"
	"github.com/ashriths/go-graph/system"
	"log"
	"math"
	"sync"
)

type MemKVStore struct {
	clock     uint64
	data      map[string]string
	listData  map[string]*list.List
	clockLock sync.Mutex
	dataLock  sync.Mutex
}

func NewMemKVStore() *MemKVStore {
	return &MemKVStore{
		data:     make(map[string]string),
		listData: make(map[string]*list.List),
	}
}

func (self *MemKVStore) Clock(atLeast uint64) uint64 {
	self.clockLock.Lock()
	defer self.clockLock.Unlock()

	if self.clock < atLeast {
		self.clock = atLeast
	}

	ret := self.clock

	if self.clock < math.MaxUint64 {
		self.clock++
	}

	if system.Logging {
		log.Printf("Clock(%d) => %d", atLeast, ret)
	}

	return ret
}

func (self *MemKVStore) Get(key string) string {
	self.dataLock.Lock()
	defer self.dataLock.Unlock()

	value := self.data[key]

	if system.Logging {
		log.Printf("Get(%q) => %q", key, value)
	}

	return value
}

func (self *MemKVStore) Set(key string, value string) error {
	self.dataLock.Lock()
	defer self.dataLock.Unlock()

	if value != "" {
		self.data[key] = value
	} else {
		delete(self.data, key)
	}

	if system.Logging {
		log.Printf("Set(%q, %q)", key, value)
	}

	return nil
}

func (self *MemKVStore) Keys(p Pattern) []string {
	self.dataLock.Lock()
	defer self.dataLock.Unlock()

	ret := make([]string, 0, len(self.data))

	for k := range self.data {
		if p.Match(k) {
			ret = append(ret, k)
		}
	}

	if system.Logging {
		log.Printf("Keys(%q, %q) => %d", p.Prefix, p.Suffix, len(ret))
		for i, s := range ret {
			log.Printf("  %d: %q", i, s)
		}
	}

	return ret
}

func (self *MemKVStore) ListGet(key string) []string {
	self.dataLock.Lock()
	defer self.dataLock.Unlock()
	var retList []string
	if lst, found := self.listData[key]; !found {
		return []string{}
	} else {
		retList = make([]string, 0, lst.Len())
		for i := lst.Front(); i != nil; i = i.Next() {
			retList = append(retList, i.Value.(string))
		}
	}

	if system.Logging {
		log.Printf("ListGet(%q) => %d", key, len(retList))
		for i, s := range retList {
			log.Printf("  %d: %q", i, s)
		}
	}

	return retList
}

func (self *MemKVStore) ListAppend(key string, value string) error {
	self.dataLock.Lock()
	defer self.dataLock.Unlock()

	lst, found := self.listData[key]
	if !found {
		lst = list.New()
		self.listData[key] = lst
	}

	lst.PushBack(value)

	if system.Logging {
		log.Printf("ListAppend(%q, %q)", key, value)
	}

	return nil
}

func (self *MemKVStore) ListRemove(key string, value string) int {
	self.dataLock.Lock()
	defer self.dataLock.Unlock()

	n := 0

	lst, found := self.listData[key]
	if !found {
		return n
	}

	i := lst.Front()
	for i != nil {
		if i.Value.(string) == value {
			hold := i
			i = i.Next()
			lst.Remove(hold)
			n++
			continue
		}

		i = i.Next()
	}

	if lst.Len() == 0 {
		delete(self.listData, key)
	}

	if system.Logging {
		log.Printf("ListRemove(%q, %q) => %d", key, value, n)
	}

	return n
}

func (self *MemKVStore) ListKeys(p Pattern) []string {
	self.dataLock.Lock()
	defer self.dataLock.Unlock()

	ret := make([]string, 0, len(self.listData))
	for k := range self.listData {
		if p.Match(k) {
			ret = append(ret, k)
		}
	}

	if system.Logging {
		log.Printf("ListKeys(%q, %q) => %d", p.Prefix, p.Suffix, len(ret))
		for i, s := range ret {
			log.Printf("  %d: %q", i, s)
		}
	}

	return ret
}

var _ KVStore = new(MemKVStore)
