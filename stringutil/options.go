package stringutil

// SlugOption configures Slug behaviour.
type SlugOption func(*slugConfig)

// slugConfig holds settings for Slug.
type slugConfig struct {
	separator string
	maxLength int
}

func defaultSlugConfig() slugConfig {
	return slugConfig{separator: "-", maxLength: 0}
}

// WithSeparator sets the character inserted between words in the slug.
// The separator must be a non-alphanumeric string (default "-").
func WithSeparator(sep string) SlugOption {
	return func(c *slugConfig) {
		if sep != "" {
			c.separator = sep
		}
	}
}

// WithMaxLength limits the slug to at most n bytes. The result is
// truncated at the nearest separator boundary to avoid mid-word cuts.
// A value of 0 (or negative) means no limit.
func WithMaxLength(n int) SlugOption {
	return func(c *slugConfig) {
		if n > 0 {
			c.maxLength = n
		}
	}
}
