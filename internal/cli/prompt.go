package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// promptLine prints label, reads one line from r, and returns the trimmed value.
func promptLine(r *bufio.Reader, w io.Writer, label string) (string, error) {
	fmt.Fprintf(w, "%s: ", boldStyle.Render(label))
	line, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("prompt: reading %q: %w", label, err)
	}
	return strings.TrimSpace(line), nil
}

// promptSecret reads a sensitive value. Input is hidden when stdin is a terminal.
func promptSecret(r *bufio.Reader, w io.Writer, label string) (string, error) {
	fmt.Fprintf(w, "%s: ", boldStyle.Render(label))

	if term.IsTerminal(int(os.Stdin.Fd())) {
		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(w) // newline after masked input
		if err != nil {
			return "", fmt.Errorf("prompt: reading %q: %w", label, err)
		}
		return strings.TrimSpace(string(b)), nil
	}

	// non-interactive: read plaintext from pipe
	line, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("prompt: reading %q: %w", label, err)
	}
	return strings.TrimSpace(line), nil
}

// confirmPrompt asks a yes/no question and returns true only for an explicit "y" or "yes".
func confirmPrompt(r *bufio.Reader, w io.Writer, question string) bool {
	fmt.Fprintf(w, "%s [y/N]: ", question)
	line, _ := r.ReadString('\n')
	v := strings.TrimSpace(strings.ToLower(line))
	return v == "y" || v == "yes"
}
