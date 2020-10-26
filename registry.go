package isvalid

import (
	"regexp"
	"sync"
)

var regexpCache = struct {
	m  map[string]*regexp.Regexp
	mu sync.Mutex
}{m: make(map[string]*regexp.Regexp)}

func RegisterRegexp(expr string) {
	regexpCache.mu.Lock()
	defer regexpCache.mu.Unlock()

	re := regexp.MustCompile(expr)
	regexpCache.m[expr] = re
}
