package cli

import (
	"fmt"

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

func runDoctor(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 2
	fmt.Fprintln(cmd.OutOrStdout(), "doctor not yet implemented")
	return nil
}
