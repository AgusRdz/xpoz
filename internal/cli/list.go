package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved tunnels",
	Args:  cobra.NoArgs,
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 1
	fmt.Fprintln(cmd.OutOrStdout(), "list not yet implemented")
	return nil
}
