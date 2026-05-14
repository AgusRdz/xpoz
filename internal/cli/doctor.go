package cli

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/AgusRdz/xpoz/internal/config"
	"github.com/AgusRdz/xpoz/internal/dns"
	"github.com/AgusRdz/xpoz/internal/store"
	"github.com/AgusRdz/xpoz/internal/tunnel"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Verify that everything is configured correctly",
	Long: `Run a series of health checks:
  - config file exists and is valid
  - cloudflared is installed and reachable
  - wildcard CNAME exists in Cloudflare DNS
  - internal proxy port is available
  - tunnel database is healthy`,
	Args: cobra.NoArgs,
	RunE: runDoctor,
}

func runDoctor(cmd *cobra.Command, _ []string) error {
	out := cmd.OutOrStdout()
	ctx := cmd.Context()

	fmt.Fprintln(out)
	fmt.Fprintln(out, boldStyle.Render("xpoz doctor"))
	fmt.Fprintln(out)

	failed := 0
	run := func(label string, fn func() (string, bool)) {
		if !printCheck(out, label, fn) {
			failed++
		}
	}

	var cfg *config.Config
	run("Config file is valid", func() (string, bool) {
		c, err := config.Load()
		if err != nil {
			return err.Error(), false
		}
		cfg = c
		return config.Path(), true
	})

	run("cloudflared binary is installed", checkCloudflared)
	run("cloudflared config exists", checkCloudflaredConfig)
	run("Proxy port is available", checkProxyPort(cfg))
	run("Store database is accessible", checkStoreDB)

	if cfg != nil {
		run("Wildcard CNAME exists in DNS", checkDNS(ctx, cfg))
	}

	fmt.Fprintln(out)
	if failed == 0 {
		printSuccess(out, "All checks passed — xpoz is ready.")
		return nil
	}
	return fmt.Errorf("doctor: %d check(s) failed", failed)
}

func printCheck(out io.Writer, label string, fn func() (string, bool)) bool {
	detail, ok := fn()
	sym := successStyle.Render("✓")
	if !ok {
		sym = errorStyle.Render("✗")
	}
	fmt.Fprintf(out, "  %s  %s\n", sym, label)
	if detail != "" {
		fmt.Fprintf(out, "     %s\n", dimStyle.Render(detail))
	}
	return ok
}

func checkCloudflared() (string, bool) {
	path, err := exec.LookPath("cloudflared")
	if err != nil {
		return "Install cloudflared: https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/downloads/", false
	}
	return path, true
}

func checkCloudflaredConfig() (string, bool) {
	p := tunnel.DefaultConfigPath()
	if _, err := os.Stat(p); err != nil {
		return "Run 'cloudflared tunnel create <name>' to set up a tunnel first", false
	}
	return p, true
}

func checkProxyPort(cfg *config.Config) func() (string, bool) {
	return func() (string, bool) {
		port := config.DefaultProxyPort
		if cfg != nil {
			port = cfg.Proxy.Port
		}
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return fmt.Sprintf("Port %d is in use — stop any running xpoz instances first", port), false
		}
		_ = l.Close()
		return fmt.Sprintf(":%d is free", port), true
	}
}

func checkStoreDB() (string, bool) {
	s, err := store.New(config.DBPath())
	if err != nil {
		return err.Error(), false
	}
	_ = s.Close()
	return config.DBPath(), true
}

func checkDNS(ctx context.Context, cfg *config.Config) func() (string, bool) {
	return func() (string, bool) {
		client, err := dns.New(cfg.Cloudflare.Token, cfg.Cloudflare.ZoneID)
		if err != nil {
			return err.Error(), false
		}
		records, err := client.ListRecords(ctx, "*."+cfg.Domain.Name)
		if err != nil {
			return err.Error(), false
		}
		if len(records) == 0 {
			return "Run 'xpoz setup' to create the wildcard CNAME record", false
		}
		return fmt.Sprintf("*.%s → %s", cfg.Domain.Name, records[0].Content), true
	}
}
