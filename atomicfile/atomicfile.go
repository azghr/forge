// Package atomicfile provides atomic file operations: writing or replacing
// files without leaving partial data on failure.
//
// It wraps OS operations (write-to-temp + rename). On POSIX, rename is
// atomic; on Windows, os.Rename is not fully atomic (see docs).
package atomicfile

import (
	"io"
	"os"
	"path/filepath"
)

// WriteFile atomically writes the contents of r to filename.
// It creates a temporary file in the same directory, copies data from r,
// fsyncs (unless opted out), and renames over filename.
// If any step fails the temp file is removed and the original is left intact.
func WriteFile(filename string, r io.Reader, opts ...Option) error {
	var cfg options
	cfg.fileMode = 0644
	for _, opt := range opts {
		opt(&cfg)
	}

	dir := filepath.Dir(filename)

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

	if _, err := io.Copy(f, r); err != nil {
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

	if err := os.Rename(tmpName, filename); err != nil {
		return &WriteError{Op: "rename", Err: err}
	}

	if err := syncDir(dir); err != nil {
		return &WriteError{Op: "sync-dir", Err: err}
	}

	remove = false
	return nil
}

// ReplaceFile atomically replaces dest with source using a rename.
// Both files must be on the same filesystem. If the rename fails dest is
// left unchanged.
func ReplaceFile(source, dest string) error {
	dir := filepath.Dir(dest)

	if err := os.Rename(source, dest); err != nil {
		return &WriteError{Op: "rename", Err: err}
	}

	if err := syncDir(dir); err != nil {
		return &WriteError{Op: "sync-dir", Err: err}
	}

	return nil
}

// syncDir fsyncs the directory to ensure the directory entry is durable on
// filesystems that cache directory metadata.
func syncDir(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	err = f.Sync()
	f.Close()
	return err
}
