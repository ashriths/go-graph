package storage

import (
	"sync"
	"math"
	"log"
)

type MemStorage struct {
	clock     uint64
	data      map[string]string
	clockLock sync.Mutex
	dataLock  sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]string),
	}
}

func (self *MemStorage) Clock(atLeast uint64) (error, uint64) {
	self.clockLock.Lock()
	defer self.clockLock.Unlock()

	if self.clock < atLeast {
		self.clock = atLeast
	}

	ret := self.clock

	if self.clock < math.MaxUint64 {
		self.clock++
	}

	if Logging {
		log.Printf("Clock(%d) => %d", atLeast, ret)
	}

	return nil, ret
}

func (self *MemStorage) Get(key string) (error, string) {
	self.dataLock.Lock()
	defer self.dataLock.Unlock()

	value := self.data[key]

	if Logging {
		log.Printf("Get(%q) => %q", key, value)
	}

	return nil, value
}

func (self *MemStorage) Set(key string, value string) error {
	self.dataLock.Lock()
	defer self.dataLock.Unlock()

	if value != "" {
		self.data[key] = value
	} else {
		delete(self.data, key)
	}

	if Logging {
		log.Printf("Set(%q, %q)", key, value)
	}

	return nil
}

func (self *MemStorage) Keys(p Pattern) (error, []string) {
	self.dataLock.Lock()
	defer self.dataLock.Unlock()

	ret := make([]string, 0, len(self.data))

	for k := range self.data {
		if p.Match(k) {
			ret = append(ret, k)
		}
	}

	if Logging {
		log.Printf("Keys(%q, %q) => %d", p.Prefix, p.Suffix, len(ret))
		for i, s := range ret {
			log.Printf("  %d: %q", i, s)
		}
	}

	return nil, ret
}

var _ Storage = new(MemStorage)
