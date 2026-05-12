// Package pathsafe safely joins paths to avoid directory traversal.
//
// Given a base directory and a relative path, pathsafe ensures the result is
// within base. It returns an error if path traversal is detected.
package pathsafe

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// SafeJoin joins base and rel ensuring the result is within base.
// Returns the cleaned absolute path or an error if the result escapes base.
//
// Performance: O(path length). Uses filepath.Abs and filepath.Clean.
func SafeJoin(base, rel string) (string, error) {
	return joinOpts(base, rel, options{})
}

// SafeJoinContext joins base and rel with context cancellation and functional
// options. It behaves identically to SafeJoin when no options are provided.
//
// The context is checked before any filesystem operations begin. If ctx is
// done, ctx.Err() is returned immediately.
func SafeJoinContext(ctx context.Context, base, rel string, opts ...Option) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return joinOpts(base, rel, o)
}

func joinOpts(base, rel string, o options) (string, error) {
	baseAbs, err := filepath.Abs(base)
	if err != nil {
		return "", fmt.Errorf("pathsafe: resolving base path: %w", err)
	}
	baseAbs = filepath.Clean(baseAbs)

	joined := filepath.Join(baseAbs, rel)
	joined = filepath.Clean(joined)

	if o.followSymlinks {
		resolved, err := filepath.EvalSymlinks(joined)
		if err != nil {
			return "", fmt.Errorf("pathsafe: resolving symlinks: %w", err)
		}
		joined = resolved

		baseResolved, err := filepath.EvalSymlinks(baseAbs)
		if err != nil {
			return "", fmt.Errorf("pathsafe: resolving base symlinks: %w", err)
		}
		baseAbs = filepath.Clean(baseResolved)
	}

	if !isWithin(baseAbs, joined) {
		return "", ErrOutsideBase
	}

	return joined, nil
}

// isWithin reports whether path is within base or equal to base.
func isWithin(base, path string) bool {
	if path == base {
		return true
	}
	return strings.HasPrefix(path, base+string(filepath.Separator))
}
