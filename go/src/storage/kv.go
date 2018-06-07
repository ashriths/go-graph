package storage

import "strings"

type Pattern struct {
	Prefix string
	Suffix string
}

type KeyValue struct {
	Key   string
	Value string
}

func KV(k, v string) *KeyValue { return &KeyValue{k, v} }

func (p *Pattern) Match(k string) bool {
	ret := strings.HasPrefix(k, p.Prefix)
	ret = ret && strings.HasSuffix(k, p.Suffix)
	return ret
}
