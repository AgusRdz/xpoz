package cli

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Bold(true)
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")).Bold(true)
	boldStyle    = lipgloss.NewStyle().Bold(true)
	dimStyle     = lipgloss.NewStyle().Faint(true)
	urlStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#60a5fa"))
)

func printSuccess(w io.Writer, msg string) {
	fmt.Fprintln(w, successStyle.Render("✓")+" "+msg)
}

func printError(w io.Writer, msg string) {
	fmt.Fprintln(w, errorStyle.Render("✗")+" "+msg)
}

func printWarn(w io.Writer, msg string) {
	fmt.Fprintln(w, warnStyle.Render("⚠")+" "+msg)
}

func printStep(w io.Writer, n int, msg string) {
	fmt.Fprintf(w, "\n%s %s\n", dimStyle.Render(fmt.Sprintf("[%d]", n)), msg)
}
