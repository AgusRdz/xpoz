package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/AgusRdz/xpoz/internal/config"
	"github.com/AgusRdz/xpoz/internal/store"
)

var startCmd = &cobra.Command{
	Use:   "start [name]",
	Short: "Start a saved tunnel in the background, or all of them",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runStart,
}

func runStart(cmd *cobra.Command, args []string) error {
	s, err := store.New(config.DBPath())
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}
	defer s.Close()

	ctx := cmd.Context()
	out := cmd.OutOrStdout()

	if len(args) == 1 {
		return startOne(ctx, out, s, args[0])
	}
	return startAll(ctx, out, s)
}

func startOne(ctx context.Context, out io.Writer, s store.Store, name string) error {
	tun, err := s.Get(ctx, name)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("no tunnel named %q — run 'xpoz list' to see saved tunnels", name)
		}
		return fmt.Errorf("looking up %q: %w", name, err)
	}

	if pidAlive(name) {
		printWarn(out, fmt.Sprintf("Tunnel %q is already running (PID file exists).", name))
		return nil
	}

	return startDetached(out, tun.LocalPort, tun.Subdomain, tun.FullURL, tun.Name)
}

func startAll(ctx context.Context, out io.Writer, s store.Store) error {
	tunnels, err := s.List(ctx)
	if err != nil {
		return fmt.Errorf("listing tunnels: %w", err)
	}

	if len(tunnels) == 0 {
		fmt.Fprintln(out, dimStyle.Render("No saved tunnels. Run 'xpoz [port] --name <name>' to create one."))
		return nil
	}

	started := 0
	for _, t := range tunnels {
		if pidAlive(t.Name) {
			fmt.Fprintf(out, dimStyle.Render("  %s already running, skipping\n"), t.Name)
			continue
		}
		if err := startDetached(out, t.LocalPort, t.Subdomain, t.FullURL, t.Name); err != nil {
			printError(out, fmt.Sprintf("%s: %s", t.Name, err.Error()))
			continue
		}
		started++
	}

	if started == 0 && len(tunnels) > 0 {
		fmt.Fprintln(out, dimStyle.Render("All tunnels are already running."))
	}
	return nil
}
