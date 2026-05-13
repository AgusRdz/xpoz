package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var flagKeepData bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove xpoz and all its data",
	Long: `Remove the xpoz binary, configuration, and tunnel database.

Use --keep-data to preserve the tunnel history database.`,
	Args: cobra.NoArgs,
	RunE: runUninstall,
}

func init() {
	uninstallCmd.Flags().BoolVar(&flagKeepData, "keep-data", false, "preserve the tunnel history database")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 2
	fmt.Fprintln(cmd.OutOrStdout(), "uninstall not yet implemented")
	return nil
}
