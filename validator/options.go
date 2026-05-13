package validator

// Option configures ValidateStruct behaviour.
type Option func(*config)

type config struct {
	tagName string
}

func defaultConfig() config {
	return config{tagName: "validate"}
}

// WithTagName sets the struct tag key used for validation rules
// (default "validate").
func WithTagName(name string) Option {
	return func(c *config) {
		if name != "" {
			c.tagName = name
		}
	}
}
