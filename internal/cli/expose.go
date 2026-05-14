package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/AgusRdz/xpoz/internal/config"
	"github.com/AgusRdz/xpoz/internal/proxy"
	"github.com/AgusRdz/xpoz/internal/store"
	"github.com/AgusRdz/xpoz/internal/subdomain"
	"github.com/spf13/cobra"
)

var (
	flagName      string
	flagReset     bool
	flagBg        bool
	flagSubdomain string // hidden: pre-resolved subdomain passed by --bg re-exec
)

func init() {
	rootCmd.Flags().StringVar(&flagName, "name", "", "assign a persistent name to the tunnel")
	rootCmd.Flags().BoolVar(&flagReset, "reset", false, "regenerate subdomain even if name already exists")
	rootCmd.Flags().BoolVar(&flagBg, "bg", false, "run in background (daemon mode)")
	rootCmd.Flags().StringVar(&flagSubdomain, "subdomain", "", "")
	_ = rootCmd.Flags().MarkHidden("subdomain")
}

func runExpose(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	port, err := strconv.Atoi(args[0])
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid port %q — must be a number between 1 and 65535", args[0])
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("not configured — run 'xpoz setup' first")
	}

	s, err := store.New(config.DBPath())
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}
	defer s.Close()

	sub, err := resolveSubdomain(cmd.Context(), s, cfg, flagName, flagReset)
	if err != nil {
		return err
	}
	hostname := sub + "." + cfg.Domain.Name
	fullURL := "https://" + hostname

	if flagBg {
		return launchBackground(cmd.OutOrStdout(), port, sub, fullURL)
	}
	return runTunnel(cmd, s, cfg, sub, hostname, fullURL, port)
}

func launchBackground(out io.Writer, port int, sub, fullURL string) error {
	args := buildChildArgs(port, sub, flagName)
	proc := exec.Command(os.Args[0], args...)
	proc.Stdin = nil
	proc.Stdout = nil
	proc.Stderr = nil
	setSysProcAttr(proc)

	if err := proc.Start(); err != nil {
		return fmt.Errorf("starting background process: %w", err)
	}

	printSuccess(out, "Tunnel started in background")
	fmt.Fprintln(out, "  "+urlStyle.Render(fullURL))
	fmt.Fprintf(out, "  %s  localhost:%d\n", dimStyle.Render("forwarding to"), port)
	fmt.Fprintf(out, "  %s  %d\n", dimStyle.Render("PID"), proc.Process.Pid)
	fmt.Fprintln(out, dimStyle.Render("  Stop with: kill "+strconv.Itoa(proc.Process.Pid)))
	fmt.Fprintln(out)
	return nil
}

func buildChildArgs(port int, sub, name string) []string {
	args := []string{strconv.Itoa(port), "--subdomain=" + sub}
	if name != "" {
		args = append(args, "--name="+name)
	}
	return args
}

func runTunnel(cmd *cobra.Command, s store.Store, cfg *config.Config, sub, hostname, fullURL string, port int) error {
	ctx := cmd.Context()
	out := cmd.OutOrStdout()

	prx := proxy.New(cfg.Proxy.Port)
	prx.AddRoute(hostname, fmt.Sprintf("localhost:%d", port))

	proxyErrCh := make(chan error, 1)
	go func() { proxyErrCh <- prx.Start() }()

	if flagName != "" {
		if err := upsertActiveTunnel(ctx, s, flagName, sub, fullURL, port); err != nil {
			return err
		}
	}

	printTunnelReady(out, fullURL, port)

	select {
	case <-ctx.Done():
	case err := <-proxyErrCh:
		if err != nil {
			return fmt.Errorf("proxy failed on port %d — is xpoz already running? (%w)", cfg.Proxy.Port, err)
		}
	}

	cleanCtx := context.Background()
	fmt.Fprintln(out)
	prx.RemoveRoute(hostname)
	_ = prx.Stop()
	if flagName != "" {
		_ = s.SetActive(cleanCtx, flagName, false)
	}
	printSuccess(out, "Tunnel closed.")
	return nil
}

// resolveSubdomain returns the subdomain for this session.
// Priority: hidden --subdomain flag (set by --bg re-exec) → named reuse → generate new.
func resolveSubdomain(ctx context.Context, s store.Store, cfg *config.Config, name string, reset bool) (string, error) {
	if flagSubdomain != "" {
		return flagSubdomain, nil
	}
	if name != "" && !reset {
		existing, err := s.Get(ctx, name)
		if err == nil {
			return existing.Subdomain, nil
		}
		if !errors.Is(err, store.ErrNotFound) {
			return "", fmt.Errorf("looking up tunnel %q: %w", name, err)
		}
	}
	taken := func(ctx context.Context, sub string) bool {
		_, err := s.GetBySubdomain(ctx, sub)
		return err == nil
	}
	sub, err := subdomain.GenerateUnique(ctx, subdomain.Style(cfg.Subdomains.Style), cfg.Subdomains.Length, taken)
	if err != nil {
		return "", fmt.Errorf("generating subdomain: %w", err)
	}
	return sub, nil
}

func upsertActiveTunnel(ctx context.Context, s store.Store, name, sub, fullURL string, port int) error {
	now := time.Now().UTC()
	return s.Upsert(ctx, &store.Tunnel{
		Name: name, Subdomain: sub, FullURL: fullURL,
		LocalPort: port, Active: true, CreatedAt: now, LastUsed: now,
	})
}

func printTunnelReady(w io.Writer, fullURL string, localPort int) {
	fmt.Fprintln(w)
	printSuccess(w, "Tunnel active")
	fmt.Fprintln(w, "  "+urlStyle.Render(fullURL))
	fmt.Fprintf(w, "  %s  localhost:%d\n", dimStyle.Render("forwarding to"), localPort)
	fmt.Fprintln(w, dimStyle.Render("  Press Ctrl+C to stop"))
	fmt.Fprintln(w)
}
