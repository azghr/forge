package envconfig

// Option configures Load behaviour.
type Option func(*config)

type config struct {
	prefix string
}

// WithPrefix prepends prefix to every environment variable name looked up
// during Load.
func WithPrefix(prefix string) Option {
	return func(c *config) {
		c.prefix = prefix
	}
}
