package multityperror

// Option configures a MultiError.
type Option func(*config)

type config struct {
	separator string
}

func defaultConfig() config {
	return config{
		separator: "; ",
	}
}

// WithSeparator sets the separator used between error messages in Error().
// The default is "; ".
func WithSeparator(sep string) Option {
	return func(c *config) {
		c.separator = sep
	}
}
