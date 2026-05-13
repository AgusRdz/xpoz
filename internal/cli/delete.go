package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a saved tunnel",
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 1
	fmt.Fprintf(cmd.OutOrStdout(), "delete %q not yet implemented\n", args[0])
	return nil
}
