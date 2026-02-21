package cmd

import (
	"fmt"
	"os"

	"github.com/mishamyrt/hapm/internal/hapm"
)

// Hapm runs command line application.
func Hapm() int {
	rootCmd := newRootCmd(os.Stdout, os.Stderr)
	rootCmd.SetArgs(os.Args[1:])

	err := rootCmd.Execute()
	if err == nil {
		return 0
	}
	if !hapm.IsHandledError(err) {
		_, _ = fmt.Fprintf(os.Stderr, "Unexpected error: %v\n", err)
	}
	return 1
}
