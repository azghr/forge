// Package cache provides a generic in-memory cache with optional TTL expiration.
//
// It solves repeated computations or lookups by storing key/value pairs.
// Entries are lazily expired on Get; no goroutines are started unless
// WithCleanupInterval is provided.
package cache

import (
	"context"
	"sync"
	"time"
)

type item[V any] struct {
	value    V
	deadline int64 // unix nano; 0 means no expiration
}

// Cache stores values of type V under keys of type K with TTL expiration.
//
// The zero value is not usable; use NewCache to create a cache.
// Cache is safe for concurrent use.
type Cache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]item[V]
	ttl   time.Duration

	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// NewCache creates a cache. If ttl > 0, values expire after ttl elapses
// from the time they are stored. Pass ttl=0 for no expiration.
func NewCache[K comparable, V any](ttl time.Duration, opts ...Option) *Cache[K, V] {
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}

	c := &Cache[K, V]{
		items:           make(map[K]item[V]),
		ttl:             ttl,
		cleanupInterval: cfg.cleanupInterval,
	}

	if c.cleanupInterval > 0 && c.ttl > 0 {
		c.stopCleanup = make(chan struct{})
		go c.cleanupLoop()
	}

	return c
}

// Set adds key to the cache with the given value. If a value already exists
// for key it is overwritten and its deadline is reset.
func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var deadline int64
	if c.ttl > 0 {
		deadline = time.Now().Add(c.ttl).UnixNano()
	}
	c.items[key] = item[V]{value: value, deadline: deadline}
}

// Get retrieves the value for key. The bool reports whether the key was
// present and not expired. Expired entries are lazily removed on Get.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	it, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}

	if c.ttl > 0 && time.Now().UnixNano() > it.deadline {
		delete(c.items, key)
		var zero V
		return zero, false
	}

	return it.value, true
}

// Delete removes key from the cache. It is a no-op if the key does not exist.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// GetOrLoad returns the cached value for key, or calls loader to produce it.
// The result is stored in the cache before returning. If the loader returns
// an error the result is not cached.
//
// ctx is passed to loader so the caller can cancel a slow load.
func (c *Cache[K, V]) GetOrLoad(ctx context.Context, key K, loader func(context.Context) (V, error)) (V, error) {
	if v, ok := c.Get(key); ok {
		return v, nil
	}

	v, err := loader(ctx)
	if err != nil {
		var zero V
		return zero, err
	}

	c.Set(key, v)
	return v, nil
}

// Len returns the number of items currently in the cache (including
// potentially expired items not yet cleaned up).
func (c *Cache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Stop stops the background cleanup goroutine if one was started via
// WithCleanupInterval. After Stop returns the caller may continue to use
// the cache normally; expired entries are removed lazily on Get.
//
// It is safe to call Stop multiple times.
func (c *Cache[K, V]) Stop() {
	if c.stopCleanup != nil {
		select {
		case <-c.stopCleanup:
		default:
			close(c.stopCleanup)
		}
	}
}

func (c *Cache[K, V]) cleanupLoop() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopCleanup:
			return
		}
	}
}

func (c *Cache[K, V]) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for k, it := range c.items {
		if c.ttl > 0 && now > it.deadline {
			delete(c.items, k)
		}
	}
}
