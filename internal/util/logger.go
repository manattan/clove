package util

import (
	"fmt"
	"os"
)

var verbose bool

// SetVerbose sets the verbose logging mode
func SetVerbose(v bool) {
	verbose = v
}

// IsVerbose returns true if verbose mode is enabled
func IsVerbose() bool {
	return verbose
}

// Verbose prints a message only if verbose mode is enabled
func Verbose(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format, args...)
		if len(format) > 0 && format[len(format)-1] != '\n' {
			fmt.Fprintln(os.Stdout)
		}
	}
}

// Info prints an informational message to stdout
func Info(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
	if len(format) > 0 && format[len(format)-1] != '\n' {
		fmt.Fprintln(os.Stdout)
	}
}
