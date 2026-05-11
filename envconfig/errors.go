package envconfig

import "fmt"

// MissingError is returned when a required environment variable is not set
// and has no default value.
type MissingError struct {
	VarName string
}

func (e *MissingError) Error() string {
	return "required env var missing: " + e.VarName
}

// ParseError is returned when an environment variable value cannot be
// converted to the target struct field type.
type ParseError struct {
	VarName string
	Value   string
	Err     error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse env var %q value %q: %v", e.VarName, e.Value, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}
