package valid

import (
	"regexp"
	"sync"
)

// RegisterRegexp compiles the given expression and caches
// the result. The given expr is assumed to be a valid regular
// expression, if it's not then RegisterRegexp will panic.
func RegisterRegexp(expr string) {
	rxcache.mu.Lock()
	defer rxcache.mu.Unlock()

	rx := regexp.MustCompile(expr)
	rxcache.m[expr] = rx
}

var rxcache = struct {
	m  map[string]*regexp.Regexp
	mu sync.Mutex
}{m: make(map[string]*regexp.Regexp)}
