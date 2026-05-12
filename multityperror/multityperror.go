package multityperror

import (
	"strings"
	"sync"
)

// MultiError aggregates multiple error values into one.
//
// It is a concurrency-safe alternative to errors.Join that supports
// incremental accumulation (append as errors occur) and custom
// formatting via functional options.
type MultiError struct {
	mu        sync.Mutex
	errs      []error
	separator string
}

// New creates a MultiError with the given options.
func New(opts ...Option) *MultiError {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return &MultiError{
		separator: cfg.separator,
	}
}

// Append adds err to the MultiError. Nil errors are silently skipped.
// It is safe for concurrent use.
func (m *MultiError) Append(err error) {
	if err == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errs = append(m.errs, err)
}

// Error implements the error interface. It returns a semicolon-separated
// message of all appended errors, or "" if no errors have been appended.
func (m *MultiError) Error() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.errs) == 0 {
		return ""
	}
	sep := m.separator
	if sep == "" {
		sep = "; "
	}
	msgs := make([]string, len(m.errs))
	for i, e := range m.errs {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, sep)
}

// IsEmpty reports whether no errors have been appended (including cases
// where only nil errors were appended).
func (m *MultiError) IsEmpty() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.errs) == 0
}

// Len returns the number of errors appended.
func (m *MultiError) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.errs)
}

// Unwrap returns the list of contained errors for compatibility with
// errors.Is and errors.As (Go 1.20+).
func (m *MultiError) Unwrap() []error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.errs) == 0 {
		return nil
	}
	out := make([]error, len(m.errs))
	copy(out, m.errs)
	return out
}

// Errors returns a copy of the contained error slice.
func (m *MultiError) Errors() []error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.errs) == 0 {
		return nil
	}
	out := make([]error, len(m.errs))
	copy(out, m.errs)
	return out
}
