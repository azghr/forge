package jsonmerge

// Option configures Merge behaviour.
type Option func(*mergeConfig)

type mergeConfig struct {
	sliceMode SliceMode
}

func defaultConfig() mergeConfig {
	return mergeConfig{sliceMode: SliceReplace}
}

// WithSliceMode sets the merge behaviour for slices.
func WithSliceMode(m SliceMode) Option {
	return func(c *mergeConfig) {
		c.sliceMode = m
	}
}
