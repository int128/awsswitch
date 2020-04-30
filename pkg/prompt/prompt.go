// Package prompt provides the user interaction features.
package prompt

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// ReadPassword reads a password from stdin with mask.
func ReadPassword(message string) (string, error) {
	if _, err := fmt.Fprint(os.Stderr, message); err != nil {
		return "", fmt.Errorf("cannot write to stderr: %w", err)
	}
	b, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("cannot read from stdin: %w", err)
	}
	if _, err := fmt.Fprintln(os.Stderr); err != nil {
		return "", fmt.Errorf("cannot write to stderr: %w", err)
	}
	return string(b), nil
}
