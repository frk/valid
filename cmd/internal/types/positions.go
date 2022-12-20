package types

import (
	"sync"
)

var positions = struct {
	sync.RWMutex
	// m should map created *StructFields to their corresponding source
	// location as obtained from search.FileAndLine.
	m map[any]string
}{m: make(map[any]string)}

func storePosition(v any, pos string) {
	positions.Lock()
	positions.m[v] = pos
	positions.Unlock()
}

func GetPosition(v any) string {
	positions.RLock()
	pos := positions.m[v]
	positions.RUnlock()
	return pos
}
