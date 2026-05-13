package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive initial configuration",
	Long:  `Configure your domain, Cloudflare API token, and Tunnel ID.`,
	Args:  cobra.NoArgs,
	RunE:  runSetup,
}

func runSetup(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 1
	fmt.Fprintln(cmd.OutOrStdout(), "setup not yet implemented")
	return nil
}
