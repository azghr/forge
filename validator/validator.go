// Package validator validates struct fields via `validate` tags.
// It supports a small set of built-in rules: nonzero and email.
//
// The package uses reflection once per call, has zero external dependencies,
// and is fully concurrency-safe.
package validator

import (
	"fmt"
	"regexp"
	"reflect"
	"strings"
)

// emailRegex is a relaxed pattern for common email addresses.
var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+$`)

// ValidateStruct validates all exported fields of v based on the configured
// validation tags. v must be a struct (or a pointer to a struct); non-struct
// values return an error.
//
// Supported tag rules (comma-separated):
//
//	nonzero — field must not be the zero value (non-empty string, non-nil
//	         pointer, non-zero number, etc.).
//	email   — string field must match a basic email pattern.
//
// Example tag: `validate:"nonzero,email"`
func ValidateStruct(v any, opts ...Option) error {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return fmt.Errorf("validator: nil pointer")
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("validator: expected struct, got %s", rv.Kind())
	}

	rt := rv.Type()
	var errs ValidationError

	for i := range rv.NumField() {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}

		tag := field.Tag.Get(cfg.tagName)
		if tag == "" {
			continue
		}

		rules := strings.Split(tag, ",")
		val := rv.Field(i)

		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if rule == "" {
				continue
			}
			switch rule {
			case "nonzero":
				if isZero(val) {
					errs = append(errs, FieldError{
						Field: field.Name,
						Tag:   "nonzero",
						Value: val.Interface(),
					})
				}
			case "email":
				if val.Kind() != reflect.String {
					errs = append(errs, FieldError{
						Field: field.Name,
						Tag:   "email",
						Value: val.Interface(),
					})
					continue
				}
				s := val.String()
				if s == "" || !emailRegex.MatchString(s) {
					errs = append(errs, FieldError{
						Field: field.Name,
						Tag:   "email",
						Value: s,
					})
				}
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// isZero reports whether v is the zero value for its type.
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.IsNil() || v.Len() == 0
	case reflect.Array:
		return v.Len() == 0
	default:
		return false
	}
}
