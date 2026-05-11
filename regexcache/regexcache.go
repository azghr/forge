// Package regexcache caches compiled regex patterns to avoid regexp.Compile
// overhead when the same patterns are used repeatedly.
package regexcache

import (
	"regexp"
	"sync"
)

// Cache is a concurrency-safe cache of compiled regex patterns.
// The zero value is not safe to use; use New to create a Cache.
type Cache struct {
	mu      sync.RWMutex
	m       map[string]*regexp.Regexp
	maxSize int
}

// New creates a new Cache with optional configuration.
func New(opts ...Option) *Cache {
	c := &Cache{m: make(map[string]*regexp.Regexp)}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Compile returns a cached *regexp.Regexp for pattern.
// If pattern has been compiled before, the cached value is returned
// immediately. Otherwise it is compiled, stored, and returned.
// If pattern is not a valid regex, an error is returned.
func (c *Cache) Compile(pattern string) (*regexp.Regexp, error) {
	c.mu.RLock()
	r, ok := c.m[pattern]
	c.mu.RUnlock()
	if ok {
		return r, nil
	}

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if r, ok := c.m[pattern]; ok {
		return r, nil
	}
	if c.maxSize > 0 && len(c.m) >= c.maxSize {
		for k := range c.m {
			delete(c.m, k)
			break
		}
	}
	c.m[pattern] = compiled
	return compiled, nil
}

// MustCompile returns a cached *regexp.Regexp for pattern.
// It panics if pattern is not a valid regex.
func (c *Cache) MustCompile(pattern string) *regexp.Regexp {
	r, err := c.Compile(pattern)
	if err != nil {
		panic("regexcache: Compile(" + pattern + "): " + err.Error())
	}
	return r
}
