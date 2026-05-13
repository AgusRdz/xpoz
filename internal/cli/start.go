package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start all persistent tunnels",
	Args:  cobra.NoArgs,
	RunE:  runStart,
}

func runStart(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 2
	fmt.Fprintln(cmd.OutOrStdout(), "start not yet implemented")
	return nil
}
