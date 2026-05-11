package atomicfile

import "os"

// Option configures the write behaviour.
type Option func(*options)

type options struct {
	fileMode os.FileMode
	noFSync  bool
}

// WithFileMode sets the file permission bits for the written file.
// The default is 0644.
func WithFileMode(mode os.FileMode) Option {
	return func(o *options) {
		o.fileMode = mode
	}
}

// WithoutFSync disables the fsync call before the rename. Disabling fsync
// improves throughput but risks data loss or corruption after an unclean
// shutdown.
func WithoutFSync() Option {
	return func(o *options) {
		o.noFSync = true
	}
}
