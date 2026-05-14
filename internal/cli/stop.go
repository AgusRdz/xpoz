package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/AgusRdz/xpoz/internal/config"
	"github.com/AgusRdz/xpoz/internal/store"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [name]",
	Short: "Stop a running background tunnel, or all of them",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runStop,
}

func runStop(cmd *cobra.Command, args []string) error {
	s, err := store.New(config.DBPath())
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}
	defer s.Close()

	ctx := cmd.Context()
	out := cmd.OutOrStdout()

	if len(args) == 1 {
		return stopOne(ctx, out, s, args[0])
	}
	return stopAll(ctx, out, s)
}

func stopOne(ctx context.Context, out io.Writer, s store.Store, name string) error {
	if err := killByName(name); err != nil {
		printWarn(out, err.Error())
	}
	if err := s.SetActive(ctx, name, false); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("no tunnel named %q — run 'xpoz list' to see saved tunnels", name)
		}
		return fmt.Errorf("updating store: %w", err)
	}
	printSuccess(out, fmt.Sprintf("Tunnel %q stopped.", name))
	return nil
}

func stopAll(ctx context.Context, out io.Writer, s store.Store) error {
	tunnels, err := s.List(ctx)
	if err != nil {
		return fmt.Errorf("listing tunnels: %w", err)
	}

	stopped := 0
	for _, t := range tunnels {
		if !t.Active && !pidAlive(t.Name) {
			continue
		}
		if err := killByName(t.Name); err != nil {
			printWarn(out, fmt.Sprintf("%s: %s", t.Name, err.Error()))
		}
		_ = s.SetActive(ctx, t.Name, false)
		printSuccess(out, fmt.Sprintf("Tunnel %q stopped.", t.Name))
		stopped++
	}

	if stopped == 0 {
		fmt.Fprintln(out, dimStyle.Render("No active tunnels."))
	}
	return nil
}

// killByName sends an interrupt to the process recorded in the PID file for name.
// The background process handles cleanup (route removal, store update) on interrupt.
func killByName(name string) error {
	pid, err := readPID(name)
	if err != nil {
		return fmt.Errorf("%q was not started with --bg or PID file is missing", name)
	}
	removePID(name)

	if !isRunning(pid) {
		return nil // already exited; PID file was stale
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil // race: process exited between the isRunning check and here
	}
	if err := proc.Signal(os.Interrupt); err != nil {
		return proc.Kill() // fallback if interrupt is unsupported
	}
	return nil
}
