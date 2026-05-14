// Package cli implements the xpoz command-line interface.
package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "xpoz [port]",
	Short: "Expose local ports via your own domain and Cloudflare Tunnel",
	Long: `xpoz exposes local ports to the internet using your own domain and Cloudflare Tunnel.

No central server. No xpoz account. Just your domain, cloudflared, and this tool.

Run 'xpoz setup' to get started.`,
	Args:          cobra.MaximumNArgs(1),
	RunE:          runExpose,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute is the CLI entry point called from main.
// It wires a signal-aware context so Ctrl+C cancels the running command.
func Execute(version string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rootCmd.Version = version
	rootCmd.SetVersionTemplate("xpoz {{.Version}}\n")
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.AddCommand(
		setupCmd,
		listCmd,
		deleteCmd,
		startCmd,
		stopCmd,
		serviceCmd,
		doctorCmd,
		updateCmd,
		autoUpdateCmd,
		uninstallCmd,
		configureCmd,
		completionCmd,
		versionCmd,
	)
}
