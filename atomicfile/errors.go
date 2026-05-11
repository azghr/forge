package atomicfile

import "fmt"

// WriteError is returned when a write operation fails.
// The Op field identifies the failing step ("create", "write", "fsync",
// "close", "rename", "sync-dir"), and Err wraps the underlying OS error.
type WriteError struct {
	Op  string
	Err error
}

func (e *WriteError) Error() string {
	return fmt.Sprintf("atomicfile: %s: %v", e.Op, e.Err)
}

func (e *WriteError) Unwrap() error {
	return e.Err
}
