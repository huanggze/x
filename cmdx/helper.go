package cmdx

import (
	"fmt"
	"io"
)

var (
	// ErrNilDependency is returned if a dependency is missing.
	ErrNilDependency = fmt.Errorf("a dependency was expected to be defined but is nil. Please open an issue with the stack trace")
	// ErrNoPrintButFail is returned to detect a failure state that was already reported to the user in some way
	ErrNoPrintButFail = fmt.Errorf("this error should never be printed")

	debugStdout, debugStderr = io.Discard, io.Discard
)
