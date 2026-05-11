package option_test

import (
	"testing"

	"github.com/azghr/forge/option"
)

func TestOption(t *testing.T) {
	t.Parallel()

	t.Run("Some", func(t *testing.T) {
		o := option.Some(5)
		if !o.IsSome() {
			t.Error("expected IsSome()=true")
		}
		v, ok := o.Unwrap()
		if !ok {
			t.Error("expected ok=true")
		}
		if v != 5 {
			t.Errorf("expected 5, got %v", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		o := option.None[int]()
		if o.IsSome() {
			t.Error("expected IsSome()=false")
		}
		_, ok := o.Unwrap()
		if ok {
			t.Error("expected ok=false")
		}
	})

	t.Run("Must panics on None", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic on Must()")
			}
		}()
		_ = option.None[string]().Must()
	})

	t.Run("Must returns value on Some", func(t *testing.T) {
		o := option.Some("hello")
		if v := o.Must(); v != "hello" {
			t.Errorf("expected 'hello', got %q", v)
		}
	})

	t.Run("zero value string", func(t *testing.T) {
		o := option.None[string]()
		v, ok := o.Unwrap()
		if ok {
			t.Error("expected ok=false")
		}
		if v != "" {
			t.Errorf("expected empty string, got %q", v)
		}
	})

	t.Run("zero value struct", func(t *testing.T) {
		type S struct{ X int }
		o := option.None[S]()
		v, ok := o.Unwrap()
		if ok {
			t.Error("expected ok=false")
		}
		if v.X != 0 {
			t.Errorf("expected zero struct, got %+v", v)
		}
	})
}

func TestOptionConcurrentSafe(t *testing.T) {
	o := option.Some(42)
	run := make(chan struct{})
	done := make(chan struct{}, 20)

	for range 20 {
		go func() {
			<-run
			o.IsSome()
			o.Unwrap()
			o.Must()
			done <- struct{}{}
		}()
	}

	close(run)
	for range 20 {
		<-done
	}
}

func BenchmarkOptionSome(b *testing.B) {
	o := option.Some(42)
	b.ResetTimer()
	for b.Loop() {
		_, _ = o.Unwrap()
	}
}

func BenchmarkOptionNone(b *testing.B) {
	o := option.None[int]()
	b.ResetTimer()
	for b.Loop() {
		_, ok := o.Unwrap()
		_ = ok
	}
}
