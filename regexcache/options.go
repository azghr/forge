package regexcache

// Option configures a Cache.
type Option func(*Cache)

// WithMaxSize limits the cache to at most n entries.
// When the cache exceeds this size, existing entries are evicted.
// A value of 0 (the default) means no limit.
func WithMaxSize(n int) Option {
	return func(c *Cache) {
		c.maxSize = n
	}
}
