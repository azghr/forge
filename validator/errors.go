package validator

import "fmt"

// FieldError describes a single field validation failure.
type FieldError struct {
	// Field is the struct field name.
	Field string
	// Tag is the validation tag that failed (e.g. "nonzero", "email").
	Tag string
	// Value is the actual value of the field.
	Value any
}

// ValidationError is a slice of FieldError values, one per failed field.
// It implements the error interface.
type ValidationError []FieldError

func (ve ValidationError) Error() string {
	switch len(ve) {
	case 0:
		return "no errors"
	case 1:
		return fmt.Sprintf("%s: %s", ve[0].Field, ve[0].Tag)
	default:
		return fmt.Sprintf("%s: %s (and %d more errors)", ve[0].Field, ve[0].Tag, len(ve)-1)
	}
}
