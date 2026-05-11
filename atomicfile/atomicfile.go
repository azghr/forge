// Package atomicfile provides atomic file writes using a temp-file + rename
// pattern. Writes are visible atomically at the target path: either the
// original content is preserved or the new content appears in full, never a
// partial or corrupted state.
//
// The typical flow is:
//  1. Create a temporary file in the same directory as the target path.
//  2. Write data to the temporary file.
//  3. Optionally fsync the temporary file (default: on).
//  4. Rename the temporary file over the target path (atomic on POSIX).
//  5. On any error, clean up the temporary file.
//
// Directories in the path other than the last element must already exist.
package atomicfile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Error sentinel values.
var (
	// ErrCancelled is returned when a context is cancelled before the
	// operation completes.
	ErrCancelled = errors.New("atomicfile: operation cancelled")
)

// WriteError is returned when the write phase (create, write, fsync, or
// rename) fails. Callers can use errors.Is to distinguish write failures
// from other error types.
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

// Write writes data to path atomically. It creates a temporary file in the
// same directory, writes data, fsyncs (unless opted out), and renames it
// over path. If any step fails the temp file is removed.
func Write(path string, data []byte, opts ...Option) error {
	return write(context.Background(), path, func(f *os.File) error {
		_, err := f.Write(data)
		return err
	}, opts)
}

// WriteContext is like Write but aborts if ctx is done before the rename.
func WriteContext(ctx context.Context, path string, data []byte, opts ...Option) error {
	return write(ctx, path, func(f *os.File) error {
		_, err := f.Write(data)
		return err
	}, opts)
}

// WriteReader copies from r into path atomically. The reader is consumed
// fully; if it returns an error the temp file is discarded.
func WriteReader(ctx context.Context, path string, r io.Reader, opts ...Option) error {
	return write(ctx, path, func(f *os.File) error {
		_, err := io.Copy(f, r)
		return err
	}, opts)
}

func write(ctx context.Context, path string, writeFn func(*os.File) error, opts []Option) error {
	if err := ctx.Err(); err != nil {
		return ErrCancelled
	}

	var cfg options
	cfg.fileMode = 0644
	for _, opt := range opts {
		opt(&cfg)
	}

	dir := filepath.Dir(path)

	f, err := os.CreateTemp(dir, ".atomic-*")
	if err != nil {
		return &WriteError{Op: "create", Err: err}
	}
	tmpName := f.Name()

	if cfg.fileMode != 0 {
		if err := os.Chmod(tmpName, cfg.fileMode); err != nil {
			f.Close()
			os.Remove(tmpName)
			return &WriteError{Op: "create", Err: err}
		}
	}

	remove := true
	defer func() {
		if remove {
			os.Remove(tmpName)
		}
	}()

	if err := writeFn(f); err != nil {
		f.Close()
		return &WriteError{Op: "write", Err: err}
	}

	if !cfg.noFSync {
		if err := f.Sync(); err != nil {
			f.Close()
			return &WriteError{Op: "fsync", Err: err}
		}
	}

	if err := f.Close(); err != nil {
		return &WriteError{Op: "close", Err: err}
	}

	if err := ctx.Err(); err != nil {
		return ErrCancelled
	}

	if err := os.Rename(tmpName, path); err != nil {
		return &WriteError{Op: "rename", Err: err}
	}

	if err := syncDir(dir); err != nil {
		return &WriteError{Op: "sync-dir", Err: err}
	}

	remove = false
	return nil
}

// syncDir fsyncs the parent directory to ensure the directory entry for the
// renamed file is durable on filesystems that cache directory metadata.
func syncDir(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	err = f.Sync()
	f.Close()
	return err
}
