package isvalid

import (
	"regexp"
	"sync"
)

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
