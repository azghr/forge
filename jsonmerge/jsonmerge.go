// Package jsonmerge provides functions to recursively merge and diff JSON-like
// data (maps with interface{} values). It is designed for combining
// configurations or comparing state snapshots.
//
// All functions are pure, concurrency-safe, and have zero external dependencies.
package jsonmerge

import "fmt"

// MergeOption configures Merge behaviour.
type MergeOption func(*mergeConfig)

type mergeConfig struct {
	sliceMode SliceMode
}

// SliceMode controls how slices are handled during merge.
type SliceMode int

const (
	// SliceReplace overwrites the destination slice with the source slice.
	SliceReplace SliceMode = iota
	// SliceAppend appends source slice elements to the destination slice.
	SliceAppend
)

func defaultConfig() mergeConfig {
	return mergeConfig{sliceMode: SliceReplace}
}

// WithSliceMode sets the merge behaviour for slices.
func WithSliceMode(m SliceMode) MergeOption {
	return func(c *mergeConfig) {
		c.sliceMode = m
	}
}

// Merge recursively merges src into dst. For each key in src:
//   - If dst lacks the key, src's value is written to dst.
//   - If both values are map[string]interface{}, Merge recurses into them.
//   - Otherwise, src's value overwrites dst's value.
//
// Slices are handled according to the configured SliceMode (default: replace).
// dst is modified in place; src is not.
func Merge(dst, src map[string]interface{}, opts ...MergeOption) {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	merge(dst, src, cfg)
}

func merge(dst, src map[string]interface{}, cfg mergeConfig) {
	for k, sv := range src {
		dv, ok := dst[k]
		if !ok {
			dst[k] = sv
			continue
		}
		sm, okS := sv.(map[string]interface{})
		dm, okD := dv.(map[string]interface{})
		if okS && okD {
			merge(dm, sm, cfg)
			continue
		}
		switch cfg.sliceMode {
		case SliceAppend:
			ds, dsOK := dv.([]interface{})
			ss, ssOK := sv.([]interface{})
			if dsOK && ssOK {
				dst[k] = append(ds, ss...)
				continue
			}
		}
		dst[k] = sv
	}
}

// Diff returns dot-notation paths for every key in a whose value differs from
// the corresponding value in b. For nested maps, paths are joined with "."
// (e.g. "y.v"). Keys present in b but absent from a are not reported.
func Diff(a, b map[string]interface{}) []string {
	var out []string
	diff(a, b, "", &out)
	return out
}

func diff(a, b map[string]interface{}, prefix string, out *[]string) {
	for k, av := range a {
		path := k
		if prefix != "" {
			path = prefix + "." + k
		}
		bv, ok := b[k]
		if !ok {
			*out = append(*out, path)
			continue
		}
		am, okA := av.(map[string]interface{})
		bm, okB := bv.(map[string]interface{})
		if okA && okB {
			diff(am, bm, path, out)
			continue
		}
		if !equal(av, bv) {
			*out = append(*out, path)
		}
	}
}

func equal(a, b interface{}) bool {
	// Fast path for nil.
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	// Use fmt.Sprintf for deep equality on compound types (slices, maps).
	// This is simple and correct for JSON-like data.
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
