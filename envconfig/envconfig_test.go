package envconfig_test

import (
	"errors"
	"os"
	"testing"

	"github.com/azghr/forge/envconfig"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("default applied when env unset", func(t *testing.T) {
		os.Clearenv()
		type C struct {
			A int `env:"A,default=5"`
		}
		var c C
		if err := envconfig.Load(&c); err != nil {
			t.Fatal(err)
		}
		if c.A != 5 {
			t.Errorf("default not applied: got %d, want 5", c.A)
		}
	})

	t.Run("required missing returns error", func(t *testing.T) {
		os.Clearenv()
		type C struct {
			B string `env:"B,required"`
		}
		var c C
		if err := envconfig.Load(&c); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("required present succeeds", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("B", "val")
		defer os.Unsetenv("B")
		type C struct {
			B string `env:"B,required"`
		}
		var c C
		if err := envconfig.Load(&c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.B != "val" {
			t.Errorf("got %q, want %q", c.B, "val")
		}
	})

	t.Run("env var overrides default", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("PORT", "9000")
		defer os.Unsetenv("PORT")
		type C struct {
			Port int `env:"PORT,default=8080"`
		}
		var c C
		if err := envconfig.Load(&c); err != nil {
			t.Fatal(err)
		}
		if c.Port != 9000 {
			t.Errorf("got %d, want 9000", c.Port)
		}
	})

	t.Run("bool true", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("F", "true")
		defer os.Unsetenv("F")
		type C struct {
			F bool `env:"F,default=false"`
		}
		var c C
		if err := envconfig.Load(&c); err != nil {
			t.Fatal(err)
		}
		if c.F != true {
			t.Error("expected true")
		}
	})

	t.Run("bool false", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("F", "false")
		defer os.Unsetenv("F")
		type C struct {
			F bool `env:"F,default=true"`
		}
		var c C
		if err := envconfig.Load(&c); err != nil {
			t.Fatal(err)
		}
		if c.F != false {
			t.Error("expected false")
		}
	})

	t.Run("float", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("R", "3.14")
		defer os.Unsetenv("R")
		type C struct {
			R float64 `env:"R,default=0.0"`
		}
		var c C
		if err := envconfig.Load(&c); err != nil {
			t.Fatal(err)
		}
		if c.R != 3.14 {
			t.Errorf("got %f, want 3.14", c.R)
		}
	})

	t.Run("uint", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("U", "42")
		defer os.Unsetenv("U")
		type C struct {
			U uint `env:"U,default=0"`
		}
		var c C
		if err := envconfig.Load(&c); err != nil {
			t.Fatal(err)
		}
		if c.U != 42 {
			t.Errorf("got %d, want 42", c.U)
		}
	})

	t.Run("string without tag is skipped", func(t *testing.T) {
		os.Clearenv()
		type C struct {
			X string
		}
		var c C
		if err := envconfig.Load(&c); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("exported field without env tag keeps zero value", func(t *testing.T) {
		os.Clearenv()
		type C struct {
			Port int `env:"PORT,default=8080"`
			Host string
		}
		var c C
		if err := envconfig.Load(&c); err != nil {
			t.Fatal(err)
		}
		if c.Port != 8080 {
			t.Errorf("got %d, want 8080", c.Port)
		}
		if c.Host != "" {
			t.Errorf("expected empty, got %q", c.Host)
		}
	})

	t.Run("non-pointer error", func(t *testing.T) {
		type C struct{}
		var c C
		if err := envconfig.Load(c); err == nil {
			t.Fatal("expected error for non-pointer")
		}
	})

	t.Run("nil pointer error", func(t *testing.T) {
		type C struct{}
		if err := envconfig.Load((*C)(nil)); err == nil {
			t.Fatal("expected error for nil")
		}
	})

	t.Run("non-struct pointer error", func(t *testing.T) {
		var x int
		if err := envconfig.Load(&x); err == nil {
			t.Fatal("expected error for non-struct")
		}
	})

	t.Run("parse error on invalid int", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("N", "not_a_number")
		defer os.Unsetenv("N")
		type C struct {
			N int `env:"N"`
		}
		var c C
		err := envconfig.Load(&c)
		if err == nil {
			t.Fatal("expected parse error")
		}
		var pe *envconfig.ParseError
		if !errors.As(err, &pe) {
			t.Errorf("expected *ParseError, got %T", err)
		}
	})

	t.Run("prefix prepended to var name", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("MYAPP_PORT", "3000")
		defer os.Unsetenv("MYAPP_PORT")
		type C struct {
			Port int `env:"PORT,default=8080"`
		}
		var c C
		if err := envconfig.Load(&c, envconfig.WithPrefix("MYAPP_")); err != nil {
			t.Fatal(err)
		}
		if c.Port != 3000 {
			t.Errorf("got %d, want 3000", c.Port)
		}
	})

	t.Run("missing required includes var name", func(t *testing.T) {
		os.Clearenv()
		type C struct {
			Key string `env:"KEY,required"`
		}
		var c C
		err := envconfig.Load(&c)
		if err == nil {
			t.Fatal("expected error")
		}
		var me *envconfig.MissingError
		if !errors.As(err, &me) {
			t.Errorf("expected *MissingError, got %T", err)
		} else if me.VarName != "KEY" {
			t.Errorf("expected VarName=KEY, got %q", me.VarName)
		}
	})
}
