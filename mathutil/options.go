package mathutil

// Option configures floating-point comparison behaviour.
type Option func(*config)

// config holds settings for ApproxEqual.
type config struct {
	epsilon float64
}

func defaultConfig() config {
	return config{epsilon: 1e-9}
}

// WithEpsilon sets the tolerance for ApproxEqual. Values ≤ 0 are clamped to
// the smallest positive float64.
func WithEpsilon(eps float64) Option {
	if eps <= 0 {
		eps = 1e-300
	}
	return func(c *config) {
		c.epsilon = eps
	}
}
