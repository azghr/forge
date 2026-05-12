package pathsafe

// Option configures SafeJoinContext behavior.
type Option func(*options)

type options struct {
	followSymlinks bool
}

// AllowSymlinkFollow enables symlink resolution before the safety check.
// When set, both the base and joined paths are resolved via filepath.EvalSymlinks
// before the containment check. This prevents symlink-based traversal attacks.
//
// Use this option when the base directory or the relative path may contain
// symbolic links that could redirect outside the intended base.
func AllowSymlinkFollow() Option {
	return func(o *options) {
		o.followSymlinks = true
	}
}
