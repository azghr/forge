package pathsafe_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/azghr/forge/pathsafe"
)

func TestSafeJoin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		base    string
		rel     string
		wantErr bool
		wantPre string
	}{
		{
			name:    "simple subpath",
			base:    "/tmp/base",
			rel:     "sub/val",
			wantErr: false,
			wantPre: "/tmp/base",
		},
		{
			name:    "dot rel",
			base:    "/tmp/base",
			rel:     ".",
			wantErr: false,
			wantPre: "/tmp/base",
		},
		{
			name:    "empty rel",
			base:    "/tmp/base",
			rel:     "",
			wantErr: false,
			wantPre: "/tmp/base",
		},
		{
			name:    "file in base",
			base:    "/tmp/base",
			rel:     "file.txt",
			wantErr: false,
			wantPre: "/tmp/base",
		},
		{
			name:    "nested subpath",
			base:    "/tmp/base",
			rel:     "a/b/c/d",
			wantErr: false,
			wantPre: "/tmp/base",
		},
		{
			name:    "traversal simple",
			base:    "/tmp/base",
			rel:     "../etc/passwd",
			wantErr: true,
		},
		{
			name:    "traversal deep",
			base:    "/tmp/base",
			rel:     "sub/../../etc",
			wantErr: true,
		},
		{
			name:    "traversal dotdot",
			base:    "/tmp/base",
			rel:     "..",
			wantErr: true,
		},
		{
			name:    "base with trailing slash",
			base:    "/tmp/base/",
			rel:     "sub/val",
			wantErr: false,
			wantPre: "/tmp/base",
		},
		{
			name:    "similar prefix not a match",
			base:    "/tmp/base",
			rel:     "../base-other",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := pathsafe.SafeJoin(tt.base, tt.rel)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SafeJoin(%q, %q) expected error, got path %q", tt.base, tt.rel, p)
				}
				if !errors.Is(err, pathsafe.ErrOutsideBase) {
					t.Errorf("SafeJoin(%q, %q) expected ErrOutsideBase, got %v", tt.base, tt.rel, err)
				}
				return
			}
			if err != nil {
				t.Errorf("SafeJoin(%q, %q) unexpected error: %v", tt.base, tt.rel, err)
			}
			if p == "" {
				t.Errorf("SafeJoin(%q, %q) returned empty path", tt.base, tt.rel)
			}
			if tt.wantPre != "" && !hasPrefix(p, tt.wantPre) {
				t.Errorf("SafeJoin(%q, %q) = %q, want path starting with %q", tt.base, tt.rel, p, tt.wantPre)
			}
		})
	}
}

func TestSafeJoinFollowSymlinks(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	baseDir := filepath.Join(dir, "base")
	targetDir := filepath.Join(dir, "target")

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a symlink inside base that points outside.
	linkPath := filepath.Join(baseDir, "escape")
	if err := os.Symlink(targetDir, linkPath); err != nil {
		t.Skip("symlink creation not supported:", err)
	}

	t.Run("symlink traversal blocked when follow enabled", func(t *testing.T) {
		_, err := pathsafe.SafeJoin(
			baseDir, "escape",
			pathsafe.AllowSymlinkFollow(),
		)
		if !errors.Is(err, pathsafe.ErrOutsideBase) {
			t.Errorf("expected ErrOutsideBase for symlink escape, got %v", err)
		}
	})

	t.Run("symlink traversal allowed without follow option", func(t *testing.T) {
		p, err := pathsafe.SafeJoin(baseDir, "escape")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !hasPrefix(p, baseDir) {
			t.Errorf("expected prefix %q, got %q", baseDir, p)
		}
	})
}

func TestSafeJoinConcurrent(t *testing.T) {
	errs := make(chan error, 20)

	for range 20 {
		go func() {
			_, err := pathsafe.SafeJoin("/tmp/base", "sub/val")
			errs <- err
		}()
	}

	for range 20 {
		if err := <-errs; err != nil {
			t.Errorf("concurrent SafeJoin failed: %v", err)
		}
	}
}

func BenchmarkSafeJoin(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		pathsafe.SafeJoin("/tmp/base", "sub/dir/file.txt")
	}
}

func BenchmarkSafeJoinTraversal(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		pathsafe.SafeJoin("/tmp/base", "../etc/passwd")
	}
}

func BenchmarkSafeJoinWithOption(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		pathsafe.SafeJoin("/tmp/base", "sub/dir/file.txt", pathsafe.AllowSymlinkFollow())
	}
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
