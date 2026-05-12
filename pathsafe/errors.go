package pathsafe

import "errors"

// ErrOutsideBase is returned when the joined path escapes the base directory.
var ErrOutsideBase = errors.New("pathsafe: result path is outside base directory")
