package isvalid

import (
	"regexp"
	"sync"
)

// RegisterRegexp compiles the given expression and caches the result. The given expr is
// assumed to be a valid regular expression, if it's not then RegisterRegexp will panic.
func RegisterRegexp(expr string) {
	regexpCache.mu.Lock()
	defer regexpCache.mu.Unlock()

	re := regexp.MustCompile(expr)
	regexpCache.m[expr] = re
}

var regexpCache = struct {
	m  map[string]*regexp.Regexp
	mu sync.Mutex
}{m: make(map[string]*regexp.Regexp)}
