package cmdx

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// ErrNilDependency is returned if a dependency is missing.
	ErrNilDependency = fmt.Errorf("a dependency was expected to be defined but is nil. Please open an issue with the stack trace")
	// ErrNoPrintButFail is returned to detect a failure state that was already reported to the user in some way
	ErrNoPrintButFail = fmt.Errorf("this error should never be printed")

	debugStdout, debugStderr = io.Discard, io.Discard
)

func init() {
	if os.Getenv("DEBUG") != "" {
		debugStdout = os.Stdout
		debugStderr = os.Stderr
	}
}

// FailSilently is supposed to be used within a commands RunE function.
// It silences cobras error handling and returns the ErrNoPrintButFail error.
func FailSilently(cmd *cobra.Command) error {
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	return errors.WithStack(ErrNoPrintButFail)
}

// Must fatals with the optional message if err is not nil.
// Deprecated: do not use this function in commands, as it makes it impossible to test them. Instead, return the error.
func Must(err error, message string, args ...interface{}) {
	if err == nil {
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, message+"\n", args...)
	os.Exit(1)
}

// Fatalf prints to os.Stderr and exists with code 1.
// Deprecated: do not use this function in commands, as it makes it impossible to test them. Instead, return the error.
func Fatalf(message string, args ...interface{}) {
	if len(args) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, message+"\n", args...)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, message)
	}
	os.Exit(1)
}
