package validator_test

import (
	"testing"

	"github.com/azghr/forge/validator"
)

func TestValidateStruct(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		type T struct {
			Name  string `validate:"nonzero"`
			Email string `validate:"email"`
		}
		u := T{Name: "Alice", Email: "x@x"}
		err := validator.ValidateStruct(u)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("nonzero failure", func(t *testing.T) {
		type T struct {
			A string `validate:"nonzero"`
			B string `validate:"nonzero"`
		}
		t1 := T{A: "ok", B: ""}
		err := validator.ValidateStruct(t1)
		if err == nil {
			t.Fatal("expected error")
		}
		ve, ok := err.(validator.ValidationError)
		if !ok {
			t.Fatalf("expected ValidationError, got %T", err)
		}
		if len(ve) != 1 {
			t.Fatalf("expected 1 error, got %d: %v", len(ve), ve)
		}
		if ve[0].Field != "B" || ve[0].Tag != "nonzero" {
			t.Errorf("unexpected error: %+v", ve[0])
		}
	})

	t.Run("multiple failures", func(t *testing.T) {
		type T struct {
			A string `validate:"nonzero"`
			B string `validate:"nonzero"`
		}
		t1 := T{A: "", B: ""}
		err := validator.ValidateStruct(t1)
		if err == nil {
			t.Fatal("expected error")
		}
		ve, ok := err.(validator.ValidationError)
		if !ok {
			t.Fatalf("expected ValidationError, got %T", err)
		}
		if len(ve) != 2 {
			t.Fatalf("expected 2 errors, got %d: %v", len(ve), ve)
		}
	})

	t.Run("email failure", func(t *testing.T) {
		type T struct {
			Email string `validate:"email"`
		}
		t1 := T{Email: "not-an-email"}
		err := validator.ValidateStruct(t1)
		if err == nil {
			t.Fatal("expected error")
		}
		ve, ok := err.(validator.ValidationError)
		if !ok {
			t.Fatalf("expected ValidationError, got %T", err)
		}
		if len(ve) != 1 || ve[0].Tag != "email" {
			t.Errorf("expected email error, got %+v", ve)
		}
	})

	t.Run("email empty", func(t *testing.T) {
		type T struct {
			Email string `validate:"email"`
		}
		t1 := T{Email: ""}
		err := validator.ValidateStruct(t1)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("email on non-string", func(t *testing.T) {
		type T struct {
			Val int `validate:"email"`
		}
		t1 := T{Val: 42}
		err := validator.ValidateStruct(t1)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("no tags", func(t *testing.T) {
		type T struct {
			A string
			B int
		}
		err := validator.ValidateStruct(T{A: "", B: 0})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("pointer to struct", func(t *testing.T) {
		type T struct {
			Name string `validate:"nonzero"`
		}
		u := &T{Name: "Alice"}
		err := validator.ValidateStruct(u)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		type T struct {
			Name string `validate:"nonzero"`
		}
		var u *T
		err := validator.ValidateStruct(u)
		if err == nil {
			t.Fatal("expected error for nil pointer")
		}
	})

	t.Run("non-struct", func(t *testing.T) {
		err := validator.ValidateStruct(42)
		if err == nil {
			t.Fatal("expected error for non-struct")
		}
	})

	t.Run("custom tag name", func(t *testing.T) {
		type T struct {
			Name string `custom:"nonzero"`
		}
		err := validator.ValidateStruct(T{Name: ""}, validator.WithTagName("custom"))
		if err == nil {
			t.Fatal("expected error with custom tag")
		}
	})

	t.Run("int nonzero", func(t *testing.T) {
		type T struct {
			Val int `validate:"nonzero"`
		}
		tests := []struct {
			name string
			val  int
			fail bool
		}{
			{name: "zero", val: 0, fail: true},
			{name: "non-zero", val: 1, fail: false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.ValidateStruct(T{Val: tt.val})
				if tt.fail && err == nil {
					t.Error("expected error for zero value")
				}
				if !tt.fail && err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			})
		}
	})

	t.Run("bool nonzero", func(t *testing.T) {
		type T struct {
			Val bool `validate:"nonzero"`
		}
		if err := validator.ValidateStruct(T{Val: false}); err == nil {
			t.Error("expected error for false bool")
		}
		if err := validator.ValidateStruct(T{Val: true}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("slice nonzero", func(t *testing.T) {
		type T struct {
			Items []int `validate:"nonzero"`
		}
		if err := validator.ValidateStruct(T{Items: nil}); err == nil {
			t.Error("expected error for nil slice")
		}
		if err := validator.ValidateStruct(T{Items: []int{}}); err == nil {
			t.Error("expected error for empty slice")
		}
		if err := validator.ValidateStruct(T{Items: []int{1}}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("pointer nonzero", func(t *testing.T) {
		type T struct {
			P *int `validate:"nonzero"`
		}
		if err := validator.ValidateStruct(T{P: nil}); err == nil {
			t.Error("expected error for nil pointer")
		}
		v := 42
		if err := validator.ValidateStruct(T{P: &v}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("email valid addresses", func(t *testing.T) {
		type T struct {
			Email string `validate:"email"`
		}
		emails := []string{
			"a@b.co",
			"user@example.com",
			"user.name+tag@example.org",
			"x@y.z",
		}
		for _, e := range emails {
			if err := validator.ValidateStruct(T{Email: e}); err != nil {
				t.Errorf("expected valid email %q, got error: %v", e, err)
			}
		}
	})

	t.Run("email invalid addresses", func(t *testing.T) {
		type T struct {
			Email string `validate:"email"`
		}
		emails := []string{
			"",
			"not-an-email",
			"@domain.com",
			"user@",
			"a b@c.com",
		}
		for _, e := range emails {
			if err := validator.ValidateStruct(T{Email: e}); err == nil {
				t.Errorf("expected invalid email %q, got no error", e)
			}
		}
	})
}

func TestConcurrentSafety(t *testing.T) {
	t.Parallel()

	type T struct {
		Name  string `validate:"nonzero"`
		Email string `validate:"email"`
		Age   int    `validate:"nonzero"`
	}

	run := make(chan struct{})
	done := make(chan struct{}, 30)

	for range 30 {
		go func() {
			<-run
			validator.ValidateStruct(T{Name: "Alice", Email: "x@x", Age: 30})
			validator.ValidateStruct(T{Name: "", Email: "bad", Age: 0})
			validator.ValidateStruct(T{Name: "Bob", Email: "b@b.com", Age: 25},
				validator.WithTagName("validate"))
			done <- struct{}{}
		}()
	}

	close(run)

	for range 30 {
		<-done
	}
}

func BenchmarkValidateValid(b *testing.B) {
	type T struct {
		Name  string `validate:"nonzero"`
		Email string `validate:"email"`
		Age   int    `validate:"nonzero"`
	}
	u := T{Name: "Alice", Email: "a@b.com", Age: 30}

	b.ResetTimer()
	for b.Loop() {
		validator.ValidateStruct(u)
	}
}

func BenchmarkValidateInvalid(b *testing.B) {
	type T struct {
		Name  string `validate:"nonzero"`
		Email string `validate:"email"`
		Age   int    `validate:"nonzero"`
	}
	u := T{Name: "", Email: "bad", Age: 0}

	b.ResetTimer()
	for b.Loop() {
		validator.ValidateStruct(u)
	}
}
