package multityperror_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/azghr/forge/multityperror"
)

type sentinelErr struct{ msg string }

func (e sentinelErr) Error() string { return e.msg }

type wrapErr struct {
	msg string
	err error
}

func (e *wrapErr) Error() string { return e.msg }
func (e *wrapErr) Unwrap() error { return e.err }

var ErrOops = sentinelErr{"oops"}

func TestMultiErrorBasic(t *testing.T) {
	t.Parallel()

	var me multityperror.MultiError
	if me.Error() != "" {
		t.Error("empty MultiError should have empty string")
	}

	me.Append(fmt.Errorf("a"))
	me.Append(fmt.Errorf("b"))

	msg := me.Error()
	if msg != "a; b" {
		t.Errorf("expected 'a; b', got %q", msg)
	}
}

func TestMultiErrorNilSkipped(t *testing.T) {
	t.Parallel()

	var me multityperror.MultiError
	me.Append(nil)
	me.Append(fmt.Errorf("err"))

	if me.IsEmpty() {
		t.Error("expected non-empty after appending non-nil error")
	}

	if n := me.Len(); n != 1 {
		t.Errorf("expected 1 error, got %d", n)
	}

	if me.Error() != "err" {
		t.Errorf("expected 'err', got %q", me.Error())
	}
}

func TestMultiErrorAllNil(t *testing.T) {
	t.Parallel()

	var me multityperror.MultiError
	me.Append(nil)
	me.Append(nil)

	if !me.IsEmpty() {
		t.Error("expected empty when all nil appended")
	}
}

func TestMultiErrorIsEmpty(t *testing.T) {
	t.Parallel()

	var me multityperror.MultiError
	if !me.IsEmpty() {
		t.Error("new MultiError should be empty")
	}

	me.Append(fmt.Errorf("x"))
	if me.IsEmpty() {
		t.Error("non-empty MultiError should not be empty")
	}
}

func TestMultiErrorLen(t *testing.T) {
	t.Parallel()

	var me multityperror.MultiError
	if n := me.Len(); n != 0 {
		t.Errorf("expected 0, got %d", n)
	}

	me.Append(fmt.Errorf("a"))
	me.Append(fmt.Errorf("b"))

	if n := me.Len(); n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}

func TestMultiErrorCustomSeparator(t *testing.T) {
	t.Parallel()

	me := multityperror.New(multityperror.WithSeparator("\n"))
	me.Append(fmt.Errorf("line1"))
	me.Append(fmt.Errorf("line2"))

	msg := me.Error()
	if msg != "line1\nline2" {
		t.Errorf("expected 'line1\\nline2', got %q", msg)
	}
}

func TestMultiErrorUnwrap(t *testing.T) {
	t.Parallel()

	var me multityperror.MultiError
	me.Append(fmt.Errorf("first"))
	me.Append(fmt.Errorf("second"))

	errs := me.Unwrap()
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
	if errs[0].Error() != "first" || errs[1].Error() != "second" {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestMultiErrorErrors(t *testing.T) {
	t.Parallel()

	var me multityperror.MultiError
	me.Append(fmt.Errorf("x"))
	me.Append(fmt.Errorf("y"))

	errs := me.Errors()
	if len(errs) != 2 {
		t.Fatalf("expected 2, got %d", len(errs))
	}

	// modify returned slice to ensure it's a copy
	errs[0] = nil
	if me.Len() != 2 {
		t.Error("modifying returned slice should not affect MultiError")
	}
}

func TestMultiErrorIs(t *testing.T) {
	t.Parallel()

	var me multityperror.MultiError
	me.Append(fmt.Errorf("some error"))
	me.Append(ErrOops)

	if !errors.Is(&me, ErrOops) {
		t.Error("errors.Is should find sentinel error in MultiError")
	}
}

func TestMultiErrorIsNotFound(t *testing.T) {
	t.Parallel()

	var me multityperror.MultiError
	me.Append(fmt.Errorf("something else"))

	if errors.Is(&me, ErrOops) {
		t.Error("errors.Is should NOT find non-existent sentinel")
	}
}

func TestMultiErrorAs(t *testing.T) {
	t.Parallel()

	target := &wrapErr{msg: "wrapped", err: fmt.Errorf("inner")}

	var me multityperror.MultiError
	me.Append(target)

	var inner *wrapErr
	if !errors.As(&me, &inner) {
		t.Error("errors.As should find wrapped error in MultiError")
	}
}

func TestMultiErrorConcurrentAppend(t *testing.T) {
	var me multityperror.MultiError

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			me.Append(fmt.Errorf("err"))
		}()
	}
	wg.Wait()

	if me.Len() != 10 {
		t.Errorf("expected 10 errors, got %d", me.Len())
	}
}

func TestMultiErrorConcurrentReadWrite(t *testing.T) {
	var me multityperror.MultiError

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			me.Append(fmt.Errorf("err"))
		}()
	}

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = me.Error()
			_ = me.IsEmpty()
			me.Len()
		}()
	}
	wg.Wait()
}

func TestMultiErrorEmptyUnwrap(t *testing.T) {
	t.Parallel()

	var me multityperror.MultiError
	errs := me.Unwrap()
	if errs != nil {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestMultiErrorEmptyErrors(t *testing.T) {
	t.Parallel()

	var me multityperror.MultiError
	errs := me.Errors()
	if errs != nil {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestMultiErrorTableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		sep       string
		errs      []error
		wantEmpty bool
		wantLen   int
		wantMsg   string
	}{
		{
			name:      "no errors",
			sep:       "; ",
			errs:      nil,
			wantEmpty: true,
			wantLen:   0,
			wantMsg:   "",
		},
		{
			name:      "all nil",
			sep:       "; ",
			errs:      []error{nil, nil},
			wantEmpty: true,
			wantLen:   0,
			wantMsg:   "",
		},
		{
			name:      "single error",
			sep:       "; ",
			errs:      []error{fmt.Errorf("one")},
			wantEmpty: false,
			wantLen:   1,
			wantMsg:   "one",
		},
		{
			name:      "multiple errors with default sep",
			sep:       "; ",
			errs:      []error{fmt.Errorf("a"), fmt.Errorf("b")},
			wantEmpty: false,
			wantLen:   2,
			wantMsg:   "a; b",
		},
		{
			name:      "mix nil and non-nil",
			sep:       "; ",
			errs:      []error{nil, fmt.Errorf("x"), nil, fmt.Errorf("y")},
			wantEmpty: false,
			wantLen:   2,
			wantMsg:   "x; y",
		},
		{
			name:      "custom separator",
			sep:       " | ",
			errs:      []error{fmt.Errorf("a"), fmt.Errorf("b")},
			wantEmpty: false,
			wantLen:   2,
			wantMsg:   "a | b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := multityperror.New(multityperror.WithSeparator(tt.sep))
			for _, e := range tt.errs {
				me.Append(e)
			}

			if me.IsEmpty() != tt.wantEmpty {
				t.Errorf("IsEmpty: expected %v, got %v", tt.wantEmpty, me.IsEmpty())
			}
			if me.Len() != tt.wantLen {
				t.Errorf("Len: expected %d, got %d", tt.wantLen, me.Len())
			}
			if me.Error() != tt.wantMsg {
				t.Errorf("Error: expected %q, got %q", tt.wantMsg, me.Error())
			}
		})
	}
}

func BenchmarkAppend(b *testing.B) {
	var me multityperror.MultiError
	err := fmt.Errorf("test error")
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		me.Append(err)
	}
}

func BenchmarkError(b *testing.B) {
	var me multityperror.MultiError
	for range 10 {
		me.Append(fmt.Errorf("error"))
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = me.Error()
	}
}

func BenchmarkIs(b *testing.B) {
	var me multityperror.MultiError
	me.Append(fmt.Errorf("some error"))
	me.Append(ErrOops)
	target := ErrOops
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		errors.Is(&me, target)
	}
}
