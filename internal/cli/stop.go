package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all running tunnels",
	Args:  cobra.NoArgs,
	RunE:  runStop,
}

func runStop(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 2
	fmt.Fprintln(cmd.OutOrStdout(), "stop not yet implemented")
	return nil
}
