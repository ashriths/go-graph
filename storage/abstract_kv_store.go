package storage

import (
	"bytes"
	"strings"
)

type KVStore interface {
	Clock(atLeast uint64) uint64
	Get(key string) string
	Set(key string, value string) error
	Keys(p Pattern) []string

	ListGet(key string) []string
	ListAppend(key, value string) error
	ListRemove(key, value string) int
	ListKeys(p Pattern) []string
}

type StorageConfig struct {
	Addr  string
	Ready chan<- bool
	Store IOMapper
}

func escapeKey(s string) string {
	s = strings.Replace(s, `|`, `||`, -1) // double all backslash
	s = strings.Replace(s, `:`, `|;`, -1) // and the replace all colons
	return s
}

func unescapeKey(s string) string {
	out := new(bytes.Buffer)
	esc := false

	for _, r := range s {
		if !esc {
			if r == '|' {
				esc = true
			} else {
				out.WriteRune(r)
			}
		} else {
			// escaping
			if r == ';' {
				out.WriteRune(':')
			} else if r == '|' {
				out.WriteRune('|')
			} else {
				// should not happen
				out.WriteRune(r)
			}

			esc = false
		}
	}

	return out.String()
}
