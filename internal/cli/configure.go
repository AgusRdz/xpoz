package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "config",
	Short: "Show or edit xpoz configuration",
	Long: `Show the current xpoz configuration or open it in $EDITOR.

Subcommands:
  xpoz config          show config file path and contents
  xpoz config edit     open config in $EDITOR`,
	Args: cobra.NoArgs,
	RunE: runConfigShow,
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open config in $EDITOR",
	Args:  cobra.NoArgs,
	RunE:  runConfigEdit,
}

func init() {
	configureCmd.AddCommand(configEditCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 1
	fmt.Fprintln(cmd.OutOrStdout(), "config not yet implemented")
	return nil
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 1
	fmt.Fprintln(cmd.OutOrStdout(), "config edit not yet implemented")
	return nil
}
